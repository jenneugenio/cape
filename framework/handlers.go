package framework

import (
	"context"
	"net/http"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/version"
	"github.com/manifoldco/go-base64"
)

// VersionResponse represents the data returned when querying the version
// handler
type VersionResponse struct {
	InstanceID string `json:"instance_id"`
	Version    string `json:"version"`
	BuildDate  string `json:"build_date"`
}

// VersionHandler returns the version information for this instance of cape.
func VersionHandler(instanceID string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, &VersionResponse{
			InstanceID: instanceID,
			Version:    version.Version,
			BuildDate:  version.BuildDate,
		})
	})
}

type LoginRequest struct {
	Email   *models.Email       `json:"email"`
	TokenID *database.ID        `json:"token_id"`
	Secret  primitives.Password `json:"secret"`
}

func LoginHandler(db database.Backend, capedb db.Interface, cp auth.CredentialProducer, ta *auth.TokenAuthority) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input LoginRequest

		err := decodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, r.URL.Path, errors.Wrap(BadJSONCause, err))
			return
		}

		logger := Logger(r.Context())

		if input.TokenID != nil {
			if err := input.TokenID.Validate(); err != nil {
				respondWithError(w, r.URL.Path, err)
				return
			}
		}

		if input.Email == nil && input.TokenID == nil {
			respondWithError(w, r.URL.Path, errors.New(InvalidParametersCause, "An email or token_id must be provided"))
			return
		}
		if input.Email != nil && input.TokenID != nil {
			respondWithError(w, r.URL.Path, errors.New(InvalidParametersCause, "You can only provide an email or a token_id."))
			return
		}

		provider, err := getCredentialProvider(r.Context(), db, capedb, input)
		if err != nil {
			logger.Info().Err(err).Msgf("Could not retrieve user for create session request, email: %s token_id: %s", input.Email, input.TokenID)
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		creds, err := provider.GetCredentials()
		if err != nil {
			logger.Info().Err(err).Msg("Could not retrieve credentials")
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		err = cp.Compare(input.Secret, &models.Credentials{
			Secret: creds.Secret,
			Salt:   creds.Salt,
			Alg:    models.CredentialsAlgType(creds.Alg),
		})
		if err != nil {
			logger.Info().Err(err).Msgf("Invalid credentials provided")
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		session, err := primitives.NewSession(provider)
		if err != nil {
			logger.Info().Err(err).Msg("Could not create session")
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		token, expiresAt, err := ta.Generate(session.ID)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to generate auth token")
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		session.SetToken(token, expiresAt)
		err = db.Create(r.Context(), session)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create session in database")
			respondWithError(w, r.URL.Path, auth.ErrAuthentication)
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token.String(),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		respondWithJSON(w, http.StatusOK, session)
	})
}

type SetupRequest struct {
	Name     models.Name         `json:"name"`
	Email    models.Email        `json:"email"`
	Password primitives.Password `json:"password"`
}

func SetupHandler(db database.Backend, capedb db.Interface, cp auth.CredentialProducer, ta *auth.TokenAuthority, rootKey [32]byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := Logger(ctx)

		var input SetupRequest
		err := decodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, r.URL.Path, errors.Wrap(BadJSONCause, err))
			return
		}

		// Since we set some key state as a part of this flow, we need to roll it back in
		// the event of an error.
		cleanup := func(err error) error {
			db.SetEncryptionCodec(nil)
			ta.SetKeyPair(nil)
			return err
		}

		doWork := func() (*models.User, error) {
			// We must create the config and load up the state before we can make
			// requests against the backend that requires the encryptionKey.
			config, encryptionKey, kp, err := createConfig(rootKey)
			if err != nil {
				logger.Error().Err(err).Msg("Could not generate config")
				return nil, err
			}

			// if setup has been run we create and add the codec here
			kms, err := crypto.LoadKMS(encryptionKey)
			if err != nil {
				logger.Error().Err(err).Msg("Could not load KMS w/ Encryption Key")
				return nil, err
			}

			// XXX: Note - if you are running more than one coordinator this _will
			// not_ work. This is a big bug that we _must_ fix prior to launch.
			//
			// See: https://github.com/capeprivacy/planning/issues/1176
			db.SetEncryptionCodec(crypto.NewSecretBoxCodec(kms))
			ta.SetKeyPair(kp)

			creds, err := cp.Generate(input.Password)
			if err != nil {
				logger.Info().Err(err).Msg("Could not generate credentials")
				return nil, err
			}

			user := models.NewUser(input.Name, input.Email, creds)

			err = capedb.Users().Create(ctx, user)
			if err != nil {
				return nil, err
			}

			tx, err := db.Transaction(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Could not create transaction")
				return nil, err
			}
			defer tx.Rollback(ctx) // nolint: errcheck

			err = tx.Create(ctx, config)
			if err != nil {
				logger.Error().Err(err).Msg("Could not create config in database")
				return nil, err
			}

			err = createSystemRoles(ctx, tx)
			if err != nil {
				logger.Error().Err(err).Msg("Could not insert roles into database")
				return nil, err
			}

			err = attachDefaultPolicy(ctx, tx, capedb)
			if err != nil {
				logger.Error().Err(err).Msg("Could not attach default policies inside database")
				return nil, err
			}

			roles, err := GetRolesByLabel(ctx, tx, []primitives.Label{
				primitives.GlobalRole,
				primitives.AdminRole,
			})
			if err != nil {
				logger.Error().Err(err).Msg("Could not retrieve roles")
				return nil, err
			}

			err = CreateAssignments(ctx, tx, user.ID, roles)
			if err != nil {
				logger.Error().Err(err).Msg("Could not create assignments in database")
				return nil, err
			}

			err = tx.Commit(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Could not commit transaction")
				return nil, err
			}

			return &user, nil
		}

		user, err := doWork()
		if err != nil {
			respondWithError(w, r.URL.Path, cleanup(err))
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	})
}

type LogoutRequest struct {
	Token *base64.Value `json:"token"`
}

func LogoutHandler(backend database.Backend, ta *auth.TokenAuthority) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input LogoutRequest
		err := decodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, r.URL.Path, errors.Wrap(BadJSONCause, err))
			return
		}

		err = doLogout(ctx, backend, ta, input)
		if err != nil {
			respondWithError(w, r.URL.Path, err)
			return
		}
	})
}

func doLogout(ctx context.Context, backend database.Backend, ta *auth.TokenAuthority, input LogoutRequest) error {
	currSession := Session(ctx)
	enforcer := auth.NewEnforcer(currSession, backend)
	if input.Token == nil {
		err := enforcer.Delete(ctx, primitives.SessionType, currSession.Session.ID)
		if err != nil {
			return err
		}

		return nil
	}

	found := false
	for _, role := range currSession.Roles {
		if role.Label == primitives.AdminRole {
			found = true
		}
	}

	if !found {
		return errors.New(auth.AuthorizationFailure, "Unable to delete session")
	}

	id, err := ta.Verify(input.Token)
	if err != nil {
		return err
	}

	err = enforcer.Delete(ctx, primitives.SessionType, id)
	if err != nil {
		return err
	}

	return nil
}
