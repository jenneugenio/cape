package worker

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"time"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"

	conn "github.com/capeprivacy/cape/connector/client"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

type SchemaJobArgs struct {
	Source *coordinator.SourceResponse
}

type Worker struct {
	pool        *pgx.ConnPool
	backend     database.Backend
	config      *Config
	coordClient *coordinator.Client
	logger      *zerolog.Logger

	// Loaded late at Start
	workers *que.WorkerPool
	qc      *que.Client
}

func (w *Worker) GetSchema(j *que.Job) error {
	ctx := context.Background()
	var sja SchemaJobArgs
	err := json.Unmarshal(j.Args, &sja)
	if err != nil {
		w.logger.Err(err).Msg("Unable to unmarshall GetSchema job arguments")
		return err
	}

	service := sja.Source.Service
	cert, err := ioutil.ReadFile("connector/certs/localhost.crt")
	if err != nil {
		w.logger.Err(err).Msg("Unable to read connector cert file")
		return err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		e := errors.New(BadCertificateCause, "Bad cert for TLS")
		w.logger.Err(e).Msg("Bad cert for TLS")
		return e
	}

	connClient, err := conn.NewClient(service.Endpoint, w.coordClient.SessionToken(), certPool)
	if err != nil {
		w.logger.Err(err).Msg("Unable to create connector client")
		return err
	}

	sr, err := connClient.Schema(ctx, sja.Source.Label)
	if err != nil {
		w.logger.Err(err).Msg("Unable to get schema from connector")
		return err
	}

	schemas := sr.GetSchemas()
	schemaBlob := primitives.SchemaBlob{}

	for _, s := range schemas {
		table := map[string]string{}
		for _, field := range s.Fields {
			table[field.Name] = field.Field.String()
		}

		schemaBlob[s.Target] = table
	}

	err = w.coordClient.ReportSchema(ctx, sja.Source.ID, schemaBlob)
	if err != nil {
		w.logger.Err(err).Msg("Unable to report schema to the coordinator")
		return err
	}

	w.logger.Info().Msg(fmt.Sprintf("Reported schema for source %s", sja.Source.Label.String()))
	return nil
}

// GetSources finds all of the sources registered in the database
// It will then schedule GetSchema jobs for each of those sources.
func (w *Worker) GetSources(j *que.Job) error {
	ctx := context.Background()
	sources, err := w.coordClient.ListSources(ctx)
	if err != nil {
		w.logger.Err(err).Msg("Error listing sources")
		return err
	}

	// schedule a job to get the schema for each of the found sources
	for _, s := range sources {
		args, err := json.Marshal(&SchemaJobArgs{Source: s})
		if err != nil {
			w.logger.Err(err).Msg("Unable to decode schema job arguments")
			return err
		}

		job := &que.Job{
			Type: "GetSchema",
			Args: args,
		}

		w.logger.Info().Msg(fmt.Sprintf("Enqueueing job %s", j.Type))
		err = w.qc.Enqueue(job)
		if err != nil {
			w.logger.Err(err).Msg("Unable to enqueue GetSchema job")
			return err
		}
	}

	return nil
}

func (w *Worker) DeleteRecoveries(j *que.Job) error {
	w.logger.Info().Msg("Attempting to delete old recoveries")

	ctx := context.Background()
	recoveries, err := w.coordClient.Recoveries(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msgf("Could not retrieve a list of recoveries from the coordinator")
		return err
	}

	ids := []database.ID{}
	for _, recovery := range recoveries {
		if recovery.Expired() {
			ids = append(ids, recovery.ID)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	w.logger.Info().Msgf("Attempting to delete %d of %d recoveries", len(ids), len(recoveries))
	err = w.coordClient.DeleteRecoveries(ctx, ids)
	if err != nil {
		w.logger.Error().Err(err).Msgf("Failed to delete recoveries")
		return err
	}

	w.logger.Info().Msg("Successfully deleted recoveries")
	return nil
}

func (w *Worker) Start() error {
	ctx := context.Background()
	defer w.pool.Close()
	w.qc = que.NewClient(w.pool)

	wm := que.WorkMap{
		"GetSources":       w.GetSources,
		"GetSchema":        w.GetSchema,
		"DeleteRecoveries": w.DeleteRecoveries,
	}

	err := w.backend.Open(ctx)
	if err != nil {
		return err
	}

	workers := que.NewWorkerPool(w.qc, wm, 2)
	w.workers = workers
	go workers.Start()

	w.poll()
	return nil
}

func (w *Worker) poll() {
	for {
		jobs := []*que.Job{
			{
				Type: "GetSources",
			},
			{
				Type: "DeleteRecoveries",
			},
		}

		for _, j := range jobs {
			// Purposely ignore the error; we should recover gracefully once
			// we connect to the database again or can successfully submit a job.
			_ = w.Enqueue(j)
		}

		// XXX: With the way this is currently setup, if the amount of work
		// exceeds what workers can do then we will keep appending jobs to the
		// end of the job queue. This will result in an explosion  in the
		// number of jobs in the queue ultimately bringing everything to a
		// hault.
		time.Sleep(5 * time.Second)
	}
}

func (w *Worker) Enqueue(j *que.Job) error {
	w.logger.Info().Msgf("Enqueuing job %s", j.Type)
	err := w.qc.Enqueue(j)
	if err != nil {
		w.logger.Error().Err(err).Msgf("Could not enqueue job %s", j.Type)
		return err
	}

	w.logger.Info().Msgf("Successfully enqueued job %s", j.Type)
	return nil
}

func NewWorker(config *Config) (*Worker, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	u := config.DatabaseURL.ToURL()
	pgxcfg, err := pgx.ParseURI(u.String())
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})
	if err != nil {
		return nil, err
	}

	backend, err := database.New(u, "worker")
	if err != nil {
		return nil, err
	}

	transport, err := coordinator.NewReAuthTransport(config.CoordinatorURL, config.Token, config.Logger)
	if err != nil {
		return nil, err
	}

	return &Worker{
		pool:        pgxpool,
		backend:     backend,
		config:      config,
		logger:      config.Logger,
		coordClient: coordinator.NewClient(transport),
	}, nil
}
