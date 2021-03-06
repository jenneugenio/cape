package encrypt

import (
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
)

// CapeDBEncrypt is a postgresql implementation of the cape database interface
type CapeDBEncrypt struct {
	db    db.Interface
	codec crypto.EncryptionCodec
}

var _ db.Interface = &CapeDBEncrypt{}

func New(db db.Interface, codec crypto.EncryptionCodec) *CapeDBEncrypt {
	return &CapeDBEncrypt{
		db:    db,
		codec: codec,
	}
}

func (c *CapeDBEncrypt) Roles() db.RoleDB               { return c.db.Roles() }
func (c *CapeDBEncrypt) Users() db.UserDB               { return &userEncrypt{db: c.db.Users(), codec: c.codec} }
func (c *CapeDBEncrypt) Contributors() db.ContributorDB { return c.db.Contributors() }
func (c *CapeDBEncrypt) Projects() db.ProjectsDB        { return c.db.Projects() }
func (c *CapeDBEncrypt) Config() db.ConfigDB            { return c.db.Config() }

func (c *CapeDBEncrypt) Secrets() db.SecretDB {
	return &secretEncrypt{db: c.db.Secrets(), codec: c.codec}
}

func (c *CapeDBEncrypt) Tokens() db.TokensDB {
	return &tokensEncrypt{
		db:    c.db.Tokens(),
		codec: c.codec,
	}
}

func (c *CapeDBEncrypt) Session() db.SessionDB {
	return &sessionEncrypt{
		db:    c.db.Session(),
		codec: c.codec,
	}
}

func (c *CapeDBEncrypt) Recoveries() db.RecoveryDB {
	return &recoveriesEncrypt{
		db:    c.db.Recoveries(),
		codec: c.codec,
	}
}
