package utils

import (
	"regexp"
	"strings"
)

// Проверка валидности email
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Проверка длины строки (например, пароля или username)
func IsMinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) >= min
}
