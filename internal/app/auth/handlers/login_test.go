//nolint:all
package handlers_test

import (
	"fmt"
	"ssemu/internal/app"
	"ssemu/internal/app/auth/handlers"
	"ssemu/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginReqHandler(t *testing.T) {
	t.Run("LoginFailure when creating a new user", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth(fmt.Sprintf("%s00", test.GenerateRandomString(4)), "password")

		assert.Equal(t, handlers.LoginFailure, result)
		assert.Zero(t, uint32(client.SessionId))
	})

	t.Run("AccountError when a creating a user with an existent username", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth("testuser00", "password")

		assert.Equal(t, handlers.AccountError, result)
		assert.Zero(t, uint32(client.SessionId))
	})

	t.Run("AccountError when logging in with invalid user", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth("invalid", "pass")

		assert.Equal(t, handlers.AccountError, result)
		assert.Zero(t, uint32(client.SessionId))
	})

	t.Run("AccountBlocked when user is banned", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth("testban", "password")

		assert.Equal(t, handlers.AccountBlocked, result)
		assert.Zero(t, uint32(client.SessionId))
	})

	t.Run("Ok when sending valid user credentials", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth("testuser", "password")

		assert.Equal(t, handlers.Ok, result)
		assert.NotZero(t, uint32(client.SessionId))
	})

	t.Run("AccountError when password is wrong", func(t *testing.T) {
		client := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer client.Dispose()

		result := client.LoginAuth("testuser", "wrongpass")

		assert.Equal(t, handlers.AccountError, result)
		assert.Zero(t, client.SessionId)
	})
}
