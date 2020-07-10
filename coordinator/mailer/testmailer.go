package mailer

import (
	"context"

	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

type TestMail struct {
	To        models.Email
	Type      string
	Arguments map[string]interface{}
}

type TestMailer struct {
	Mails []*TestMail
}

func (tm *TestMailer) SendAccountRecovery(
	ctx context.Context, user *models.User, recovery *primitives.Recovery, secret primitives.Password) error {
	tm.Mails = append(tm.Mails, &TestMail{
		To:   user.Email,
		Type: "account_recovery",
		Arguments: map[string]interface{}{
			"user":     user,
			"recovery": recovery,
			"secret":   secret,
		},
	})

	return nil
}
