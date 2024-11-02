//nolint:all
package handlers_test

import (
	"ssemu/internal/app"
	"ssemu/internal/app/game/handlers"
	"ssemu/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginReqHandler(t *testing.T) {
	test.LoadResources(t)

	t.Run("NicknamePending when user has no defined nickname", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testnick", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		// act
		result := client.LoginGame(authClient.Username, authClient.SessionId)

		// assert
		assert.Equal(t, handlers.NicknamePending, result)
	})

	t.Run("Ok when user has valid nickname", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testuser", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		// act
		result := client.LoginGame(authClient.Username, authClient.SessionId)

		// assert
		assert.Equal(t, handlers.Ok, result)
	})
}
