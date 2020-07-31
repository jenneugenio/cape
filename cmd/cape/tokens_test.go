package main

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/cmd/cape/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func TestCreateToken(t *testing.T) {
	gm.RegisterTestingT(t)

	password, user := models.GenerateUser("bob", "bob@bob.bob")

	creds := primitives.GenerateCredentials()

	me := coordinator.MeResponse{User: &models.User{
		ID:    user.ID,
		Email: user.Email,
	}}

	t.Run("Can create a token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		token, err := primitives.NewToken(user.ID, creds)
		gm.Expect(err).To(gm.BeNil())

		resp := coordinator.CreateTokenResponse{
			Response: &coordinator.CreateTokenMutation{
				Secret: password,
				Token:  token,
			},
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: me,
			},
			{
				Value: resp,
			},
		})
		err = app.Run([]string{"cape", "tokens", "create"})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(3))

		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(me.User.Email.String()))

		gm.Expect(u.Calls[1].Name).To(gm.Equal("details"))
		gm.Expect(u.Calls[2].Name).To(gm.Equal("notify"))
	})

	t.Run("Can list tokens", func(t *testing.T) {
		gm.RegisterTestingT(t)

		idStrs := []string{
			"2018d9x3ntbca95dda3bu9wnrr",
			"2015338ejcum4rzncvnugucvtc",
			"2011e949qta0quff3n4yx7ny3r",
			"201dandy989092yebk2m0143p4",
		}
		ids := make([]database.ID, len(idStrs))

		for i, s := range idStrs {
			ID, err := database.DecodeFromString(s)
			gm.Expect(err).To(gm.BeNil())

			ids[i] = ID
		}

		resp := coordinator.ListTokensResponse{
			IDs: ids,
		}

		app, u := NewHarness([]*coordinator.MockResponse{
			{
				Value: me,
			},
			{
				Value: resp,
			},
		})
		err := app.Run([]string{"cape", "tokens", "list"})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(2))

		gm.Expect(u.Calls[0].Name).To(gm.Equal("table"))
		gm.Expect(u.Calls[0].Args[0]).To(gm.Equal(ui.TableHeader{"Token ID"}))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(ui.TableBody{{idStrs[0]}, {idStrs[1]}, {idStrs[2]}, {idStrs[3]}}))

		gm.Expect(u.Calls[1].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[1].Args[0]).To(gm.Equal("\nFound {{ . | toString | faded }} token{{ . | pluralize \"s\"}}\n"))
		gm.Expect(u.Calls[1].Args[1]).To(gm.Equal(4))
	})

	t.Run("Can remove a token", func(t *testing.T) {
		app, u := NewHarness([]*coordinator.MockResponse{})
		ID := "2018d9x3ntbca95dda3bu9wnrr"
		err := app.Run([]string{"cape", "tokens", "remove", ID})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(u.Calls)).To(gm.Equal(1))
		gm.Expect(u.Calls[0].Name).To(gm.Equal("template"))
		gm.Expect(u.Calls[0].Args[0]).To(gm.Equal("Removed the token with ID {{ . | toString | faded }} from Cape\n"))
		gm.Expect(u.Calls[0].Args[1]).To(gm.Equal(ID))
	})

	t.Run("Can't remove a token if you don't specify an ID", func(t *testing.T) {
		gm.RegisterTestingT(t)

		app, _ := NewHarness([]*coordinator.MockResponse{})
		err := app.Run([]string{"cape", "tokens", "remove"})
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
