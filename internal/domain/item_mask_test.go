package domain_test

import (
	"ssemu/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemMask(t *testing.T) {
	tests := []struct {
		value       uint64
		category    int
		subcategory int
		number      int
	}{
		{0, 0, 0, 0},
		{1, 0, 0, 1},
		{1000000, 1, 0, 0},
		{2020022, 2, 2, 22},
		{3070001, 3, 7, 1},
		{9131024, 9, 13, 1024},
	}
	for _, tt := range tests {
		mask := domain.ItemMask(tt.value)
		assert.Equal(t, tt.category, mask.Category())
		assert.Equal(t, tt.subcategory, mask.SubCategory())
		assert.Equal(t, tt.number, mask.Number())
	}
}
