// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/dropoutlabs/cape/database"
)

type NewUserRequest struct {
	Name string      `json:"name"`
	ID   database.ID `json:"id"`
}
