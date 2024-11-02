package domain_test

import (
	"errors"
	"ssemu/internal/domain"
	"testing"
)

func TestNewNickname(t *testing.T) {
	failures := []struct {
		input    string
		expected error
	}{
		{"n", domain.ErrValidationInvalidLength},
		{"ni", domain.ErrValidationInvalidLength},
		{"", domain.ErrValidationInvalidLength},
		{"   ", domain.ErrValidationInvalidCharacters},
		{"nickm?", domain.ErrValidationInvalidCharacters},
		{"????", domain.ErrValidationInvalidCharacters},
		{"nicknamewithover16chars", domain.ErrValidationInvalidLength},
		{"nickmé", domain.ErrValidationInvalidCharacters},
		{"nick_mé", domain.ErrValidationInvalidCharacters},
		{"______", domain.ErrValidationInvalidCharacters},
		{".nickname.", domain.ErrValidationInvalidCharacters},
		{".nick", domain.ErrValidationInvalidCharacters},
		{"nick.", domain.ErrValidationInvalidCharacters},
		{"ni.ck", domain.ErrValidationInvalidCharacters},
		{"[nic]", domain.ErrValidationInvalidCharacters},
		{"[]", domain.ErrValidationInvalidLength},
		{"..", domain.ErrValidationInvalidLength},
		{".", domain.ErrValidationInvalidLength},
		{"..nickname..", domain.ErrValidationInvalidCharacters},
		{"**nickname*", domain.ErrValidationInvalidCharacters},
		{"*nick", domain.ErrValidationInvalidCharacters},
		{"nick*", domain.ErrValidationInvalidCharacters},
		{"nick'", domain.ErrValidationInvalidCharacters},
		{"'nick'", domain.ErrValidationInvalidCharacters},
		{"nick-", domain.ErrValidationInvalidCharacters},
		{"-nick", domain.ErrValidationInvalidCharacters},
		{"-nick-", domain.ErrValidationInvalidCharacters},
		{"!nick!", domain.ErrValidationInvalidCharacters},
		{"nick~", domain.ErrValidationInvalidCharacters},
		{"nickname", nil},
		{"nicknameeeee", nil},
		{"nam", nil},
		{"nickname0001", nil},
		{"16charssssssssss", nil},
	}

	for _, tc := range failures {
		t.Run("Failure", func(t *testing.T) {
			_, err := domain.NewNickname(tc.input)
			if !errors.Is(err, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, err)
			}
		})
	}
}
