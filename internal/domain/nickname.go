package domain

import "ssemu/internal/utils"

type Nickname string

func NewNickname(value string) (Nickname, error) {
	if len(value) < 3 || len(value) > 16 {
		return "", ErrValidationInvalidLength
	}
	if !utils.StringIsAlphanumeric(value) {
		return "", ErrValidationInvalidCharacters
	}
	return Nickname(value), nil
}
