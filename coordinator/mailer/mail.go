package mailer

import (
	"context"

	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

type Mailer interface {
	SendAccountRecovery(context.Context, models.User, models.Recovery, primitives.Password) error
}
