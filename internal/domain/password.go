package domain

import (
	"ssemu/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type Password string

func NewPassword(value string) (Password, error) {
	if len(value) < 4 || len(value) > 13 {
		return "", ErrValidationInvalidLength
	}
	if !utils.StringIsAlphanumeric(value) {
		return "", ErrValidationInvalidCharacters
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Password(hashedPassword), nil
}

func (p Password) MatchPassword(plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(plainText))
	return err == nil
}
