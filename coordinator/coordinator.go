package coordinator

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NYTimes/gziphandler"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/justinas/alice"
	"github.com/manifoldco/go-base64"
	"github.com/manifoldco/healthz"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	capepg "github.com/capeprivacy/cape/coordinator/db/postgres"
	"github.com/capeprivacy/cape/coordinator/graph"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/mailer"
	"github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Coordinator is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Coordinator struct {
	cfg     *Config
	backend database.Backend
	handler http.Handler
	logger  *zerolog.Logger
	mailer  mailer.Mailer
	pool    *pgxpool.Pool

	tokenAuth          *auth.TokenAuthority
	credentialProducer auth.CredentialProducer
}

// Setup the coordinator so it's ready to be served!
func (c *Coordinator) Setup(ctx context.Context) (http.Handler, error) {
	err := c.backend.Open(ctx)
	if err != nil {
		return nil, err
	}

	cfg, err := getDatabaseConfig(ctx, c.backend)
	if err != nil {
		// if setup hasn't been run yet
		if errors.CausedBy(err, database.NotFoundCause) {
			return c.handler, nil
		}

		return nil, err
	}

	err = c.setBackendCodec(cfg)
	if err != nil {
		return nil, err
	}

	_, err = c.pool.Exec(ctx, "SELECT 1;")
	if err != nil {
		return nil, err
	}

	err = c.setTokenAuthKeyPair(cfg, c.cfg.RootKey)
	return c.handler, err
}

// Teardown the coordinator taking it back to it's start state!
func (c *Coordinator) Teardown(ctx context.Context) error {
	c.pool.Close()
	return c.backend.Close()
}

// CertFiles implements the Component interface. Coordinator doesn't support
// TLS right now so not needed!
func (c *Coordinator) CertFiles() (certFile string, keyFile string) {
	certFile = c.cfg.CertFile
	keyFile = c.cfg.KeyFile
	return
}

func (c *Coordinator) setBackendCodec(cfg *primitives.Config) error {
	// if setup has been run we create and add the codec here
	encryptionKey, err := decryptBase64s(c.cfg.RootKey, cfg.EncryptionKey)
	if err != nil {
		return err
	}

	keyURL, err := crypto.NewKeyURL(string(encryptionKey))
	if err != nil {
		return err
	}

	kms, err := crypto.LoadKMS(keyURL)
	if err != nil {
		return err
	}

	codec := crypto.NewSecretBoxCodec(kms)

	c.backend.SetEncryptionCodec(codec)

	return nil
}

func (c *Coordinator) setTokenAuthKeyPair(cfg *primitives.Config, rootKey *base64.Value) error {
	unencrypted, err := decryptBase64s(rootKey, cfg.AuthKeypair)
	if err != nil {
		return err
	}

	pkg := &auth.KeypairPackage{}
	err = json.Unmarshal(unencrypted, pkg)
	if err != nil {
		return err
	}

	kp, err := pkg.Unpackage()
	if err != nil {
		return err
	}

	c.tokenAuth.SetKeyPair(kp)

	return nil
}

// New validates the input and returns a constructed Coordinator.
//
// If the mode is set to Testing then the Coordinator will use the SHA256
// algorithm for hashing passwords. This mode should only be used within the
// context of unit & integration tests.
func New(cfg *Config, logger *zerolog.Logger, mailer mailer.Mailer) (*Coordinator, error) {
	if cfg == nil {
		return nil, errors.New(InvalidConfigCause, "Config must be provided")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	backend, err := database.New(cfg.DB.Addr.ToURL(), cfg.InstanceID.String())
	if err != nil {
		return nil, err
	}

	var rootKey [32]byte
	copy(rootKey[:], *cfg.RootKey)

	tokenAuth, err := auth.NewTokenAuthority(nil, cfg.InstanceID.String())
	if err != nil {
		return nil, err
	}

	var cp auth.CredentialProducer
	switch cfg.CredentialProducerAlg {
	case primitives.SHA256:
		cp = auth.DefaultSHA256Producer
	case primitives.Argon2ID:
		cp = auth.DefaultArgon2IDProducer
	default:
		return nil, errors.New(InvalidConfigCause, "Unknown credential producer algorithm supplied")
	}

	pgxPool := mustPgxPool(cfg.DB.Addr.ToURL().String(), cfg.InstanceID.String())

	config := &generated.Config{
		Resolvers: &graph.Resolver{
			Database:           capepg.New(pgxPool),
			Backend:            backend,
			TokenAuthority:     tokenAuth,
			RootKey:            rootKey,
			CredentialProducer: cp,
			Mailer:             mailer,
		}}

	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(*config))
	gqlHandler.SetErrorPresenter(errorPresenter)

	authenticated := framework.IsAuthenticatedMiddleware(backend, tokenAuth)

	root := http.NewServeMux()
	root.Handle("/v1", playground.Handler("GraphQL playground", "/query"))
	root.Handle("/v1/query", framework.AuthTokenMiddleware(authenticated(gqlHandler)))
	root.Handle("/v1/version", framework.VersionHandler(cfg.InstanceID.String()))
	root.Handle("/v1/login", framework.LoginHandler(backend, cp, tokenAuth))
	root.Handle("/v1/setup", framework.SetupHandler(backend, cp, tokenAuth, rootKey))
	root.Handle("/v1/logout", framework.AuthTokenMiddleware(authenticated(framework.LogoutHandler(backend, tokenAuth))))

	health := healthz.NewHandler(root)
	chain := alice.New(
		framework.RequestIDMiddleware,
		framework.LogMiddleware(logger),
		framework.RoundtripLoggerMiddleware,
		framework.RecoveryMiddleware,
		gziphandler.GzipHandler,
		cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"Accept",
				"X-Requested-With",
				"Origin",
				"Referer",
			},
			AllowCredentials: true,
		}).Handler,
	).Then(health)

	return &Coordinator{
		cfg:                cfg,
		handler:            chain,
		backend:            backend,
		logger:             logger,
		tokenAuth:          tokenAuth,
		mailer:             mailer,
		pool:               pgxPool,
		credentialProducer: cp,
	}, nil
}

func decryptBase64s(rootKey *base64.Value, data *base64.Value) ([]byte, error) {
	var key [32]byte

	copy(key[:], *rootKey)

	decrypted, err := crypto.Decrypt(key, *data)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func getDatabaseConfig(ctx context.Context, db database.Backend) (*primitives.Config, error) {
	cfg := &primitives.Config{}

	// Querying for true is weird but no easy way to query the config right now, also
	// it gets the job done.
	err := db.QueryOne(ctx, cfg, database.NewFilter(database.Where{"setup": "true"}, nil, nil))
	return cfg, err
}

func mustPgxPool(url, name string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("error parsing database config: %s", err.Error())
	}

	// Set the application name which can be used for identifying which service
	// is connecting to postgres
	cfg.ConnConfig.RuntimeParams = map[string]string{
		"application_name": name,
	}

	// Don't connect to the database until we start using the pool
	cfg.LazyConnect = true

	c, err := pgxpool.ConnectConfig(context.TODO(), cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}

	return c
}
