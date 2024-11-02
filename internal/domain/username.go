package domain

import "ssemu/internal/utils"

type Username string

func NewUsername(value string) (Username, error) {
	if len(value) < 4 || len(value) > 10 {
		return "", ErrValidationInvalidLength
	}
	if !utils.StringIsAlphanumeric(value) {
		return "", ErrValidationInvalidCharacters
	}
	return Username(value), nil
}
