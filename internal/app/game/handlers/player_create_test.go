//nolint:all
package handlers_test

import (
	"ssemu/internal/app"
	"ssemu/internal/app/game/handlers"
	"ssemu/internal/network"
	"ssemu/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlayerReqHandler(t *testing.T) {
	test.LoadResources(t)

	t.Run("NicknameAlreadyUsed if user is trying for a already used nickname", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testnick", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		client.LoginGame(authClient.Username, authClient.SessionId)

		// act
		result := sendCreatePlayerReq(client, "TestUserPlayer")

		// assert
		assert.Equal(t, handlers.NicknameAlreadyUsed, result)
	})

	t.Run("PlayerCreateSuccess if user creates the player", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testnewplr", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		loginRes := client.LoginGame(authClient.Username, authClient.SessionId)
		assert.Equal(t, handlers.NicknamePending, loginRes)

		// act
		result := sendCreatePlayerReq(client, "TestNewPlayer")

		// assert
		assert.Equal(t, handlers.PlayerCreateSuccess, result)
	})
}

func sendCreatePlayerReq(c test.TestClient, nickname string) handlers.ServerResult {
	req := network.NewPacket(handlers.CreatePlayerReq)
	defer req.Free()

	req.WriteStringSlice(nickname, 16)
	c.Send(req)

	ack := c.Receive(handlers.ServerResultAck)
	result := handlers.ServerResult(ack.ReadU32())

	return result
}
