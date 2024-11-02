package network_test

import (
	"ssemu/internal/network"
	"ssemu/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	test.LoadResources(t)
	text := "somet3xt"
	buffer := make([]byte, 256)
	pkt := network.GetPacket(buffer)
	pkt.WriteStringSlice(text, 13)
	network.Encrypt(buffer)
	network.Decrypt(buffer)
	pkt2 := network.GetPacket(buffer)
	decText := pkt2.ReadStringSlice(13)
	assert.Equal(t, text, decText)
}
