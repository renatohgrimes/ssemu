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

func TestGetNicknameAvailabilityReqHandler(t *testing.T) {
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
		result := sendGetNicknameAvailabilityReq(client, "TestUserPlayer")

		// assert
		assert.Equal(t, handlers.NicknameAlreadyUsed, result)
	})

	t.Run("NicknameAvailable if user is trying for a not used nickname", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testnick", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		client.LoginGame(authClient.Username, authClient.SessionId)

		// act
		result := sendGetNicknameAvailabilityReq(client, test.GenerateRandomString(14))

		// assert
		assert.Equal(t, handlers.NicknameAvailable, result)
	})
}

func sendGetNicknameAvailabilityReq(c test.TestClient, nickname string) handlers.ServerResult {
	req := network.NewPacket(handlers.GetNicknameAvailabilityReq)
	defer req.Free()

	req.WriteStringSlice(nickname, 16)
	c.Send(req)

	ack := c.Receive(handlers.ServerResultAck)
	result := handlers.ServerResult(ack.ReadU32())

	return result
}
