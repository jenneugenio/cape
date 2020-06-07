package mailer

import (
	"context"

	"github.com/capeprivacy/cape/primitives"
)

type Mailer interface {
	SendAccountRecovery(context.Context, *primitives.User, *primitives.Recovery, primitives.Password) error
}
