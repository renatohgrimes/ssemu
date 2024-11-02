package domain_test

import (
	"ssemu/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCharacterMask(t *testing.T) {
	tests := []struct {
		value  uint32
		gender domain.CharacterGender
		hair   byte
		face   byte
		shirt  byte
		pants  byte
	}{
		{0, domain.Male, 0, 0, 0, 0},
		{1, domain.Female, 0, 0, 0, 0},
		{2, domain.Male, 1, 0, 0, 0},
		{5, domain.Female, 2, 0, 0, 0},
		{128, domain.Male, 0, 1, 0, 0},
		{129, domain.Female, 0, 1, 0, 0},
		{8192, domain.Male, 0, 0, 1, 0},
		{16384, domain.Male, 0, 0, 2, 0},
		{8650752, domain.Male, 0, 0, 0, 1},
		{8659075, domain.Female, 1, 1, 1, 1},
		{8659077, domain.Female, 2, 1, 1, 1},
		{17301504, domain.Male, 0, 0, 0, 2},
		{17318019, domain.Female, 1, 1, 2, 2},
	}

	for _, tt := range tests {
		mask := domain.CharacterMask(tt.value)
		assert.Equal(t, tt.gender, mask.Gender())
		assert.Equal(t, tt.hair, mask.Hair())
		assert.Equal(t, tt.face, mask.Face())
		assert.Equal(t, tt.shirt, mask.Shirt())
		assert.Equal(t, tt.pants, mask.Pants())
	}
}
