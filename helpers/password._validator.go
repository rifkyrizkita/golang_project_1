package helpers

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasLowercase bool
		hasUppercase bool
		hasSymbol    bool
		hasDigit     bool
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	// Ensure at least one lowercase, one uppercase, one digit, and one symbol
	return hasLowercase && hasUppercase && hasDigit && hasSymbol
}