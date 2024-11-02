package domain

import (
	"ssemu/internal/utils"
)

var ErrValidationInvalidLength = utils.NewValidationError("invalid length")

var ErrValidationInvalidCharacters = utils.NewValidationError("invalid characters")
