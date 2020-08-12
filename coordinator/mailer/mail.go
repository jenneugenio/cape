package mailer

import (
	"context"

	"github.com/capeprivacy/cape/models"
)

type Mailer interface {
	SendAccountRecovery(context.Context, models.User, models.Recovery, models.Password) error
}
