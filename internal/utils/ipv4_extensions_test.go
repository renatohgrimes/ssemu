package utils_test

import (
	"ssemu/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPToU32(t *testing.T) {
	assert.Equal(t, uint32(2063640768), utils.IPToU32("192.168.0.123"))
	assert.Equal(t, uint32(16777343), utils.IPToU32("127.0.0.1"))
}

func TestIPToString(t *testing.T) {
	assert.Equal(t, "192.168.0.123", utils.IPToString(2063640768))
	assert.Equal(t, "127.0.0.1", utils.IPToString(16777343))
}
