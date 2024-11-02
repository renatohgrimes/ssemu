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

func TestAdminShowWindowReqHandler(t *testing.T) {
	test.LoadResources(t)

	t.Run("Allowed if user is admin", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testadmin", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		client.LoginGame(authClient.Username, authClient.SessionId)

		// act
		result := sendAdminShowWindowReq(client)

		// assert
		assert.Equal(t, handlers.Allowed, result)
	})

	t.Run("NotAllowed if user not admin", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testnotadm", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		client.LoginGame(authClient.Username, authClient.SessionId)

		// act
		result := sendAdminShowWindowReq(client)

		// assert
		assert.Equal(t, handlers.NotAllowed, result)
	})
}

func sendAdminShowWindowReq(c test.TestClient) handlers.AdminShowWindowResult {
	req := network.NewPacket(handlers.AdminShowWindowReq)
	defer req.Free()

	c.Send(req)

	ack := c.Receive(handlers.AdminShowWindowAck)
	result := handlers.AdminShowWindowResult(ack.ReadU8())

	return result
}
