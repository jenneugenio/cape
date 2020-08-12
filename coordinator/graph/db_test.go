package graph

import "github.com/capeprivacy/cape/coordinator/db"

type testDatabase struct {
	tokensDB *tokensDB
	usersDB  *usersDB
	rolesDB  *rolesDB
}

func (t testDatabase) Roles() db.RoleDB               { return t.rolesDB }
func (t testDatabase) Users() db.UserDB               { return t.usersDB }
func (t testDatabase) Projects() db.ProjectsDB        { panic("implement me") }
func (t testDatabase) Contributors() db.ContributorDB { panic("implement me") }
func (t testDatabase) Config() db.ConfigDB            { panic("implement me") }
func (t testDatabase) Secrets() db.SecretDB           { panic("implement me") }
func (t testDatabase) Session() db.SessionDB          { panic("implement me") }
func (t testDatabase) Recoveries() db.RecoveryDB      { panic("implement me") }
func (t testDatabase) Tokens() db.TokensDB            { return t.tokensDB }
