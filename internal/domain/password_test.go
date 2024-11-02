package domain_test

import (
	"ssemu/internal/domain"
	"testing"
)

func TestMatchPassword(t *testing.T) {
	plainText := "password"
	hashedPassword, err := domain.NewPassword(plainText)
	if err != nil {
		t.Error(err)
	}
	samePass := hashedPassword.MatchPassword(plainText)
	if !samePass {
		t.Error("Match password failed")
	}
	samePass = hashedPassword.MatchPassword("wrongpassword")
	if samePass {
		t.Error("Match password failed")
	}
}
