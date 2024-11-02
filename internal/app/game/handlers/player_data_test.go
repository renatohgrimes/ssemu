package handlers_test

import (
	"ssemu/internal/app"
	"ssemu/internal/app/game/handlers"
	"ssemu/internal/network"
	"ssemu/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerDataService(t *testing.T) {
	test.LoadResources(t)

	t.Run("Test Receive Player Data", func(t *testing.T) {
		// arrange
		authClient := test.NewTestClient(t, "tcp", app.AuthServerSettings.Port)
		defer authClient.Dispose()

		authClient.LoginAuth("testplrsvc", "password")

		client := test.NewTestClient(t, "tcp", app.GameServerSettings.Port)
		defer client.Dispose()

		req := network.NewPacket(handlers.LoginReq)
		defer req.Free()

		req.WriteStringSlice(authClient.Username, 13)
		req.Skip(30)
		req.WriteU32(uint32(authClient.SessionId))

		client.Send(req)
		client.Receive(handlers.LoginAck)

		// act + assert

		licenseDataAck := client.Receive(handlers.LicenseDataAck)
		assert.Equal(t, byte(2), licenseDataAck.ReadU8())
		assert.Equal(t, byte(101), licenseDataAck.ReadU8())
		assert.Equal(t, byte(102), licenseDataAck.ReadU8())

		charSlotDataAck := client.Receive(handlers.CharacterSlotDataAck)
		assert.Equal(t, byte(2), charSlotDataAck.ReadU8())
		assert.Equal(t, byte(3), charSlotDataAck.ReadU8())
		assert.Equal(t, byte(1), charSlotDataAck.ReadU8())

		client.Receive(handlers.CharacterDataAck)
		client.Receive(handlers.CharacterEquipDataAck)

		client.Receive(handlers.CharacterDataAck)
		client.Receive(handlers.CharacterEquipDataAck)

		client.Receive(handlers.InventoryDataAck)

		serverAck := client.Receive(handlers.ServerResultAck)
		assert.Equal(t, uint32(handlers.InventorySuccess), serverAck.ReadU32())

		client.Receive(handlers.StatsDataAck)

		serverAck = client.Receive(handlers.ServerResultAck)
		assert.Equal(t, uint32(handlers.LoginSuccess), serverAck.ReadU32())
	})
}
