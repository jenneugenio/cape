package coordinator

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NYTimes/gziphandler"
	"github.com/justinas/alice"
	"github.com/manifoldco/go-base64"
	"github.com/manifoldco/healthz"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/graph"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
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

	// if setup has been run we create and add the codec here
	encryptionKey, err := decryptEncryptionKey(c.cfg.RootKey, cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	kms, err := crypto.LoadKMS(encryptionKey)
	if err != nil {
		return nil, err
	}

	codec := crypto.NewSecretBoxCodec(kms)

	c.backend.SetEncryptionCodec(codec)

	return c.handler, nil
}

// Teardown the coordinator taking it back to it's start state!
func (c *Coordinator) Teardown(ctx context.Context) error {
	return c.backend.Close()
}

// CertFiles implements the Component interface. Coordinator doesn't support
// TLS right now so not needed!
func (c *Coordinator) CertFiles() (certFile string, keyFile string) {
	return
}

// New validates the input and returns a constructed Coordinator
func New(cfg *Config, logger *zerolog.Logger) (*Coordinator, error) {
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

	keypair, err := cfg.Auth.Unpackage()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := auth.NewTokenAuthority(keypair, cfg.InstanceID.String())
	if err != nil {
		return nil, err
	}

	var rootKey [32]byte
	copy(rootKey[:], *cfg.RootKey)

	config := generated.Config{Resolvers: &graph.Resolver{
		Backend:        backend,
		TokenAuthority: tokenAuth,
		RootKey:        rootKey,
	}}

	config.Directives.IsAuthenticated = framework.IsAuthenticatedDirective(backend, tokenAuth)
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	gqlHandler.SetErrorPresenter(errorPresenter)

	root := http.NewServeMux()
	root.Handle("/v1", playground.Handler("GraphQL playground", "/query"))
	root.Handle("/v1/query", framework.AuthTokenMiddleware(gqlHandler))
	root.Handle("/v1/version", framework.VersionHandler(cfg.InstanceID.String()))

	health := healthz.NewHandler(root)
	chain := alice.New(
		framework.RequestIDMiddleware,
		framework.LogMiddleware(logger),
		framework.RoundtripLoggerMiddleware,
		framework.RecoveryMiddleware,
		gziphandler.GzipHandler,
		cors.Default().Handler,
	).Then(health)

	return &Coordinator{
		cfg:     cfg,
		handler: chain,
		backend: backend,
		logger:  logger,
	}, nil
}

func decryptEncryptionKey(rootKey *base64.Value,
	encryptionKey *base64.Value) (*crypto.KeyURL, error) {
	var key [32]byte

	copy(key[:], *rootKey)

	decrypted, err := crypto.Decrypt(key, *encryptionKey)
	if err != nil {
		return nil, err
	}

	return crypto.NewKeyURL(string(decrypted))
}

func getDatabaseConfig(ctx context.Context, db database.Backend) (*primitives.Config, error) {
	cfg := &primitives.Config{}

	// Querying for true is weird but no easy way to query the config right now, also
	// it gets the job done.
	err := db.QueryOne(ctx, cfg, database.NewFilter(database.Where{"setup": "true"}, nil, nil))
	return cfg, err
}
