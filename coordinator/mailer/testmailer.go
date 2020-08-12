package mailer

import (
	"context"

	"github.com/capeprivacy/cape/models"
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
	ctx context.Context, user models.User, recovery models.Recovery, secret models.Password) error {
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
