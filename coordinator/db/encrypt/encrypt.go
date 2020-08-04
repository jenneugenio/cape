package encrypt

import (
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
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
func (c *CapeDBEncrypt) Projects() db.ProjectsDB {
	return &projectEncrypt{db: c.db.Projects(), codec: c.codec}
}
func (c *CapeDBEncrypt) Config() db.ConfigDB { return c.db.Config() }
