package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringIsAlphanumeric(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC", true},
		{"123", true},
		{"abc-123", false},
		{"$@#", false},
		{"", false},
		{"Hello123World", true},
		{"áéíóú", false},
		{"你好世界", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := StringIsAlphanumeric(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
