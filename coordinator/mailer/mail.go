package mailer

import (
	"context"

	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

type Mailer interface {
	SendAccountRecovery(context.Context, *models.User, *primitives.Recovery, primitives.Password) error
}
