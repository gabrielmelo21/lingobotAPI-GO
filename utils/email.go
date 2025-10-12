package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail verifica se o email é válido
func ValidateEmail(email string) bool {
	// Regex básico para validação de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.TrimSpace(email))
}
