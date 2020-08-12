package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/coordinator/db/encrypt"
	capepg "github.com/capeprivacy/cape/coordinator/db/postgres"
	"github.com/capeprivacy/cape/coordinator/graph"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/mailer"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Coordinator is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Coordinator struct {
	cfg     *Config
	handler http.Handler
	logger  *zerolog.Logger
	mailer  mailer.Mailer
	pool    *pgxpool.Pool
	db      db.Interface

	tokenAuth          *auth.TokenAuthority
	credentialProducer auth.CredentialProducer
}

// Setup the coordinator so it's ready to be served!
func (c *Coordinator) Setup(ctx context.Context) (http.Handler, error) {
	return c.handler, nil
}

// Teardown the coordinator taking it back to it's start state!
func (c *Coordinator) Teardown(ctx context.Context) error {
	c.pool.Close()
	return nil
}

// CertFiles implements the Component interface. Coordinator doesn't support
// TLS right now so not needed!
func (c *Coordinator) CertFiles() (certFile string, keyFile string) {
	certFile = c.cfg.CertFile
	keyFile = c.cfg.KeyFile
	return
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

	var cp auth.CredentialProducer
	switch cfg.CredentialProducerAlg {
	case models.SHA256:
		cp = auth.DefaultSHA256Producer
	case models.Argon2ID:
		cp = auth.DefaultArgon2IDProducer
	default:
		return nil, errors.New(InvalidConfigCause, "Unknown credential producer algorithm supplied")
	}

	pgxPool := mustPgxPool(cfg.DB.Addr.ToURL().String(), cfg.InstanceID.String())
	capedb := capepg.New(pgxPool)

	coor := &Coordinator{
		cfg:                cfg,
		logger:             logger,
		mailer:             mailer,
		pool:               pgxPool,
		credentialProducer: cp,
	}

	err := coor.doSetup(context.TODO(), capedb)
	if err != nil {
		return nil, err
	}

	config := generated.Config{
		Resolvers: &graph.Resolver{
			Database:           coor.db,
			CredentialProducer: cp,
			Mailer:             mailer,
		}}

	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	gqlHandler.SetErrorPresenter(errorPresenter)

	authenticated := IsAuthenticatedMiddleware(coor)

	root := http.NewServeMux()
	root.Handle("/v1", playground.Handler("GraphQL playground", "/query"))
	root.Handle("/v1/query", AuthTokenMiddleware(authenticated(gqlHandler)))
	root.Handle("/v1/version", VersionHandler(cfg.InstanceID.String()))
	root.Handle("/v1/login", LoginHandler(coor))
	root.Handle("/v1/logout", AuthTokenMiddleware(authenticated(LogoutHandler(coor))))

	health := healthz.NewHandler(root)
	chain := alice.New(
		RequestIDMiddleware,
		LogMiddleware(logger),
		RoundtripLoggerMiddleware,
		RecoveryMiddleware,
		gziphandler.GzipHandler,
	)

	if cfg.Cors.Enable {
		logger.Info().Msg("enabling CORS")
		allowOrigin := cfg.Cors.AllowOrigin
		if allowOrigin == nil {
			allowOrigin = []string{"*"}
		}
		cors := cors.New(cors.Options{
			AllowedOrigins: allowOrigin,
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"Accept",
				"X-Requested-With",
				"Origin",
				"Referer",
			},
			AllowCredentials: true,
		}).Handler

		chain = chain.Extend(alice.New(cors))
	} else {
		logger.Info().Msg("not enabling CORS")
	}

	coor.handler = chain.Then(health)

	return coor, nil
}

func (c *Coordinator) doSetup(ctx context.Context, capedb db.Interface) error {
	_, codec, kp, err := getDatabaseConfig(ctx, capedb, c.cfg.RootKey)
	if err == nil {
		ta, err := auth.NewTokenAuthority(kp, c.cfg.InstanceID.String())
		if err != nil {
			return err
		}
		c.tokenAuth = ta

		capedb := encrypt.New(capedb, codec)
		c.db = capedb

		return nil
	}

	if err.Error() != db.ErrNoRows.Error() {
		return err
	}

	if c.cfg.User == nil {
		return fmt.Errorf("user must be specified when starting coordinator for first time")
	}

	// We must create the config and load up the state before we can make
	// requests against the backend that requires the encryptionKey.
	config, encryptionKey, kp, err := createDatabaseConfig(c.cfg.RootKey)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not generate config")
		return err
	}

	// if setup has been run we create and add the codec here
	kms, err := crypto.LoadKMS(encryptionKey)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not load KMS w/ Encryption Key")
		return err
	}

	codec = crypto.NewSecretBoxCodec(kms)
	ta, err := auth.NewTokenAuthority(kp, c.cfg.InstanceID.String())
	if err != nil {
		return err
	}
	c.tokenAuth = ta

	enc := encrypt.New(capedb, codec)
	c.db = enc

	creds, err := c.credentialProducer.Generate(models.Password(c.cfg.User.Password))
	if err != nil {
		c.logger.Info().Err(err).Msg("Could not generate credentials")
		return err
	}

	user := models.NewUser(c.cfg.User.Name, c.cfg.User.Email, *creds)
	err = enc.Users().Create(ctx, user)
	if err != nil {
		return err
	}

	err = enc.Config().Create(ctx, *config)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not create config in database")
		return err
	}

	err = enc.Roles().CreateSystemRoles(ctx)
	if err != nil {
		return err
	}

	_, err = enc.Roles().SetOrgRole(ctx, c.cfg.User.Email, models.AdminRole)
	if err != nil {
		return err
	}

	return nil
}

func getDatabaseConfig(ctx context.Context, db db.Interface, rootKey string) (*models.Config, crypto.EncryptionCodec, *auth.Keypair, error) {
	cfg, err := db.Config().Get(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	u, err := crypto.NewKeyURL(rootKey)
	if err != nil {
		return nil, nil, nil, err
	}

	tmpKms, err := crypto.LoadKMS(u)
	if err != nil {
		return nil, nil, nil, err
	}

	unencrypted, err := tmpKms.Decrypt(context.TODO(), *cfg.AuthKeypair)
	if err != nil {
		return nil, nil, nil, err
	}

	kp := &auth.Keypair{}
	err = json.Unmarshal(unencrypted, kp)
	if err != nil {
		return nil, nil, nil, err
	}

	// if setup has been run we create and add the codec here
	encryptionKey, err := tmpKms.Decrypt(context.TODO(), *cfg.EncryptionKey)
	if err != nil {
		return nil, nil, nil, err
	}

	keyURL, err := crypto.NewKeyURL(string(encryptionKey))
	if err != nil {
		return nil, nil, nil, err
	}

	kms, err := crypto.LoadKMS(keyURL)
	if err != nil {
		return nil, nil, nil, err
	}

	codec := crypto.NewSecretBoxCodec(kms)

	return cfg, codec, kp, err
}

func createDatabaseConfig(rootKey string) (*models.Config, *crypto.KeyURL, *auth.Keypair, error) {
	u, err := crypto.NewKeyURL(rootKey)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptionKey := u
	if u.Type() == crypto.Base64Key {
		key, err := crypto.NewBase64KeyURL(nil)
		if err != nil {
			return nil, nil, nil, err
		}

		encryptionKey = key
	}

	kms, err := crypto.LoadKMS(u)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedKey, err := kms.Encrypt(context.TODO(), []byte(encryptionKey.String()))
	if err != nil {
		return nil, nil, nil, err
	}

	keypair, err := auth.NewKeypair()
	if err != nil {
		return nil, nil, nil, err
	}

	by, err := json.Marshal(keypair)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedAuth, err := kms.Encrypt(context.TODO(), by)
	if err != nil {
		return nil, nil, nil, err
	}

	config, err := models.NewConfig(base64.New(encryptedKey), base64.New(encryptedAuth))
	if err != nil {
		return nil, nil, nil, err
	}

	return config, encryptionKey, keypair, nil
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
