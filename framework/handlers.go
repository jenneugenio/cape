package framework

import (
	"net/http"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/version"
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
	Email   *primitives.Email   `json:"email"`
	TokenID *database.ID        `json:"token_id"`
	Secret  primitives.Password `json:"secret"`
}

func LoginHandler(db database.Backend, cp auth.CredentialProducer, ta *auth.TokenAuthority) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input LoginRequest

		err := decodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, errors.Wrap(BadJSONCause, err))
			return
		}

		logger := Logger(r.Context())

		if input.Email != nil {
			if err := input.Email.Validate(); err != nil {
				respondWithError(w, err)
				return
			}
		}

		if input.TokenID != nil {
			if err := input.TokenID.Validate(); err != nil {
				respondWithError(w, err)
				return
			}
		}

		if input.Email == nil && input.TokenID == nil {
			respondWithError(w, errors.New(InvalidParametersCause, "An email or token_id must be provided"))
			return
		}
		if input.Email != nil && input.TokenID != nil {
			respondWithError(w, errors.New(InvalidParametersCause, "You can only provide an email or a token_id."))
			return
		}

		provider, err := getCredentialProvider(r.Context(), db, input)
		if err != nil {
			logger.Info().Err(err).Msgf("Could not retrieve user for create session request, email: %s token_id: %s", input.Email, input.TokenID)
			respondWithError(w, auth.ErrAuthentication)
			return
		}

		creds, err := provider.GetCredentials()
		if err != nil {
			logger.Info().Err(err).Msg("Could not retrieve credential provider")
			respondWithError(w, err)
			return
		}

		err = cp.Compare(input.Secret, creds)
		if err != nil {
			logger.Info().Err(err).Msgf("Invalid credentials provided")
			respondWithError(w, auth.ErrAuthentication)
			return
		}

		session, err := primitives.NewSession(provider)
		if err != nil {
			logger.Info().Err(err).Msg("Could not create session")
			respondWithError(w, auth.ErrAuthentication)
			return
		}

		token, expiresAt, err := ta.Generate(session.ID)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to generate auth token")
			respondWithError(w, auth.ErrAuthentication)
			return
		}

		session.SetToken(token, expiresAt)
		err = db.Create(r.Context(), session)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create session in database")
			respondWithError(w, auth.ErrAuthentication)
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
	Name     primitives.Name     `json:"name"`
	Email    primitives.Email    `json:"email"`
	Password primitives.Password `json:"password"`
}

func SetupHandler(db database.Backend, cp auth.CredentialProducer, ta *auth.TokenAuthority, rootKey [32]byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := Logger(ctx)

		var input SetupRequest
		err := decodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, err)
			return
		}

		// Since we set some key state as a part of this flow, we need to roll it back in
		// the event of an error.
		cleanup := func(err error) error {
			db.SetEncryptionCodec(nil)
			ta.SetKeyPair(nil)
			return err
		}

		doWork := func() (*primitives.User, error) {
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

			creds, err := cp.Generate(input.Password)
			if err != nil {
				logger.Info().Err(err).Msg("Could not generate credentials")
				return nil, err
			}

			user, err := primitives.NewUser(input.Name, input.Email, creds)
			if err != nil {
				logger.Info().Err(err).Msg("Could not create user")
				return nil, err
			}

			err = tx.Create(ctx, user)
			if err != nil {
				logger.Error().Err(err).Msg("Could not insert user into database")
				return nil, err
			}

			err = createSystemRoles(ctx, tx)
			if err != nil {
				logger.Error().Err(err).Msg("Could not insert roles into database")
				return nil, err
			}

			err = attachDefaultPolicy(ctx, tx)
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

			err = CreateAssignments(ctx, tx, user, roles)
			if err != nil {
				logger.Error().Err(err).Msg("Could not create assignments in database")
				return nil, err
			}

			err = tx.Commit(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Could not commit transaction")
				return nil, err
			}

			return user, nil
		}

		user, err := doWork()
		if err != nil {
			respondWithError(w, cleanup(err))
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	})
}
