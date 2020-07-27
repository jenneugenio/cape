package coordinator

import (
	"context"
	"net/http"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/db"
	fw "github.com/capeprivacy/cape/framework"
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

func LoginHandler(coordinator *Coordinator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input LoginRequest

		db := coordinator.backend
		capedb := coordinator.db
		cp := coordinator.credentialProducer
		ta := coordinator.tokenAuth

		err := fw.DecodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, r.URL.Path, errors.Wrap(fw.BadJSONCause, err))
			return
		}

		logger := fw.Logger(r.Context())

		if input.TokenID != nil {
			if err := input.TokenID.Validate(); err != nil {
				respondWithError(w, r.URL.Path, err)
				return
			}
		}

		if input.Email == nil && input.TokenID == nil {
			respondWithError(w, r.URL.Path, errors.New(fw.InvalidParametersCause, "An email or token_id must be provided"))
			return
		}
		if input.Email != nil && input.TokenID != nil {
			respondWithError(w, r.URL.Path, errors.New(fw.InvalidParametersCause, "You can only provide an email or a token_id."))
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
			Secure:   false,
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

type LogoutRequest struct {
	Token *base64.Value `json:"token"`
}

func LogoutHandler(coordinator *Coordinator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		backend := coordinator.backend
		ta := coordinator.tokenAuth

		var input LogoutRequest
		err := fw.DecodeJSONBody(w, r, &input)
		if err != nil {
			respondWithError(w, r.URL.Path, errors.Wrap(fw.BadJSONCause, err))
			return
		}

		err = doLogout(ctx, backend, ta, input)
		if err != nil {
			respondWithError(w, r.URL.Path, err)
			return
		}
	}
}

func doLogout(ctx context.Context, backend database.Backend, ta *auth.TokenAuthority, input LogoutRequest) error {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, backend)
	if input.Token == nil {
		err := enforcer.Delete(ctx, primitives.SessionType, currSession.Session.ID)
		if err != nil {
			return err
		}

		return nil
	}

	found := false
	if currSession.Roles.Global.Label == models.AdminRole {
		found = true
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

func getCredentialProvider(ctx context.Context, q database.Querier, capedb db.Interface, input LoginRequest) (primitives.CredentialProvider, error) {
	if input.Email != nil {
		return capedb.Users().Get(ctx, *input.Email)
	}

	token := &primitives.Token{}
	err := q.Get(ctx, *input.TokenID, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
