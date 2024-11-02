package utils

import "fmt"

func NewValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}
