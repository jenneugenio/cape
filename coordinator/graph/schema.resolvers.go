package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) Setup(ctx context.Context, input model.SetupRequest) (*primitives.User, error) {
	logger := fw.Logger(ctx)

	// Since we set some key state as a part of this flow, we need to roll it back in
	// the event of an error.
	cleanup := func(err error) error {
		r.Backend.SetEncryptionCodec(nil)
		r.TokenAuthority.SetKeyPair(nil)
		return err
	}

	doWork := func() (*primitives.User, error) {
		// We must create the config and load up the state before we can make
		// requests against the backend that requires the encryptionKey.
		config, encryptionKey, kp, err := createConfig(r.RootKey)
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
		r.Backend.SetEncryptionCodec(crypto.NewSecretBoxCodec(kms))
		r.TokenAuthority.SetKeyPair(kp)

		tx, err := r.Backend.Transaction(ctx)
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

		creds, err := r.CredentialProducer.Generate(input.Password)
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

		roles, err := createSystemRoles(ctx, r.Database)
		if err != nil {
			logger.Error().Err(err).Msg("Could not insert roles into database")
			return nil, err
		}
		fmt.Println(roles)

		err = attachDefaultPolicy(ctx, tx, r.Database)
		if err != nil {
			logger.Error().Err(err).Msg("Could not attach default policies inside database")
			return nil, err
		}

		err = createAssignments(ctx, tx, user, roles)
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
		return nil, cleanup(err)
	}

	return user, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserRequest) (*model.CreateUserResponse, error) {
	logger := fw.Logger(ctx)
	session := fw.Session(ctx)

	password, err := primitives.GeneratePassword()
	if err != nil {
		logger.Error().Err(err).Msg("Could not create password")
		return nil, err
	}

	creds, err := r.CredentialProducer.Generate(password)
	if err != nil {
		logger.Info().Err(err).Msg("Could not generate credentials")
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer := auth.NewEnforcer(session, tx)

	user, err := primitives.NewUser(input.Name, input.Email, creds)
	if err != nil {
		return nil, err
	}

	// We need to get the system roles back from the database so we can
	// assignment them to this user appropriately.
	systemRole, err := queryRoleByLabel(ctx, r.Database, primitives.GlobalRole)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	err = createAssignments(ctx, tx, user, []*primitives.Role{systemRole})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{
		Password: password,
		User:     user,
	}, nil
}

func (r *mutationResolver) CreateSession(ctx context.Context, input model.SessionRequest) (*primitives.Session, error) {
	logger := fw.Logger(ctx)

	if input.Email != nil {
		if err := input.Email.Validate(); err != nil {
			return nil, err
		}
	}

	if input.TokenID != nil {
		if err := input.TokenID.Validate(); err != nil {
			return nil, err
		}
	}

	if input.Email == nil && input.TokenID == nil {
		return nil, errs.New(InvalidParametersCause, "An email or token_id must be provided")
	}
	if input.Email != nil && input.TokenID != nil {
		return nil, errs.New(InvalidParametersCause, "You can only provide an email or a token_id.")
	}

	provider, err := getCredentialProvider(ctx, r.Backend, input)
	if err != nil {
		logger.Info().Err(err).Msgf("Could not retrieve identity for create session request, email: %s token_id: %s", input.Email, input.TokenID)
		return nil, auth.ErrAuthentication
	}

	creds, err := provider.GetCredentials()
	if err != nil {
		logger.Info().Err(err).Msg("Could not retrieve credential provider")
		return nil, err
	}

	err = r.CredentialProducer.Compare(input.Secret, creds)
	if err != nil {
		logger.Info().Err(err).Msgf("Invalid credentials provided")
		return nil, auth.ErrAuthentication
	}

	session, err := primitives.NewSession(provider)
	if err != nil {
		logger.Info().Err(err).Msg("Could not create session")
		return nil, auth.ErrAuthentication
	}

	token, expiresAt, err := r.TokenAuthority.Generate(session.ID)
	if err != nil {
		logger.Info().Err(err).Msg("Failed to generate auth token")
		return nil, auth.ErrAuthentication
	}

	session.SetToken(token, expiresAt)
	err = r.Backend.Create(ctx, session)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create session in database")
		return nil, auth.ErrAuthentication
	}

	return session, nil
}

func (r *mutationResolver) DeleteSession(ctx context.Context, input model.DeleteSessionRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if input.Token == nil {
		err := enforcer.Delete(ctx, primitives.SessionType, currSession.Session.ID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	found := false
	for _, role := range currSession.Roles {
		if role.Label == primitives.AdminRole {
			found = true
		}
	}

	if !found {
		return nil, errs.New(auth.AuthorizationFailure, "Unable to delete session")
	}

	id, err := r.TokenAuthority.Verify(input.Token)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.SessionType, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) User(ctx context.Context, id database.ID) (*primitives.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	user := &primitives.User{}
	err := enforcer.Get(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*primitives.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var users []*primitives.User
	err := enforcer.Query(ctx, &users, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *queryResolver) Me(ctx context.Context) (primitives.Identity, error) {
	session := fw.Session(ctx)
	return session.Identity, nil
}

func (r *queryResolver) Identities(ctx context.Context, emails []*primitives.Email) ([]primitives.Identity, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	serviceEmails := database.In{}
	userEmails := database.In{}

	for _, email := range emails {
		if email.Type == primitives.ServiceEmail {
			serviceEmails = append(serviceEmails, email.String())
		} else {
			userEmails = append(userEmails, email.String())
		}
	}

	var users []*primitives.User
	if len(userEmails) > 0 {
		err := enforcer.Query(ctx, &users, database.NewFilter(database.Where{"email": userEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	var services []*primitives.Service
	if len(serviceEmails) > 0 {
		err := enforcer.Query(ctx, &services, database.NewFilter(database.Where{"email": serviceEmails}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	identities := make([]primitives.Identity, len(users)+len(services))
	for i, user := range users {
		identities[i] = user
	}

	for i, service := range services {
		identities[i+len(users)] = service
	}

	return identities, nil
}

func (r *userResolver) Roles(ctx context.Context, obj *primitives.User) ([]*primitives.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return fw.QueryRoles(ctx, enforcer, r.Database, obj.ID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
