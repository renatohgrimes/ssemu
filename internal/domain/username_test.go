package domain_test

import (
	"errors"
	"ssemu/internal/domain"
	"testing"
)

func TestNewUsername(t *testing.T) {
	tt := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "Valid username",
			input:       "john123",
			expectedErr: nil,
		},
		{
			name:        "Invalid length (too short)",
			input:       "a",
			expectedErr: domain.ErrValidationInvalidLength,
		},
		{
			name:        "Invalid length (too long)",
			input:       "thisisaverylongusername",
			expectedErr: domain.ErrValidationInvalidLength,
		},
		{
			name:        "Invalid characters",
			input:       "user@name",
			expectedErr: domain.ErrValidationInvalidCharacters,
		},
		{
			name:        "Empty string",
			input:       "",
			expectedErr: domain.ErrValidationInvalidLength,
		},
		{
			name:        "Too short username",
			input:       "u",
			expectedErr: domain.ErrValidationInvalidLength,
		},
		{
			name:        "Non-ASCII character username",
			input:       "úse",
			expectedErr: domain.ErrValidationInvalidCharacters,
		},
		{
			name:        "Special characters username",
			input:       "????",
			expectedErr: domain.ErrValidationInvalidCharacters,
		},
		{
			name:        "Valid username",
			input:       "user",
			expectedErr: nil,
		},
		{
			name:        "Valid username with digits",
			input:       "use123",
			expectedErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := domain.NewUsername(tc.input); !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
