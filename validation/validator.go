package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

type ValidationError struct {
	Error error
	Field string
}

var (
	isUsernameValid = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isFullNameValid = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(val string, minLength, maxLength int) error {
	length := len(val)
	if length > maxLength || length < minLength {
		return fmt.Errorf("should be of %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidatePassword(password string) *ValidationError {
	return &ValidationError{ValidateString(password, 6, 72), "password"}
}

func ValidateUsername(username string) *ValidationError {
	if err := ValidateString(username, 3, 100); err != nil {
		return &ValidationError{err, "username"}
	}
	if !isUsernameValid(username) {
		return &ValidationError{fmt.Errorf("must contain lowercase letters, digits and _ only"), "username"}
	}
	return nil
}

func ValidateFullName(username string) *ValidationError {
	if err := ValidateString(username, 3, 100); err != nil {
		return &ValidationError{err, "full_name"}
	}
	if !isFullNameValid(username) {
		return &ValidationError{fmt.Errorf("must contain letters and spaces only"), "full_name"}
	}
	return nil
}

func ValidateEmail(email string) *ValidationError {
	if err := ValidateString(email, 3, 100); err != nil {
		return &ValidationError{err, "email"}
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return &ValidationError{fmt.Errorf("is not a valid email"), "email"}
	}
	return nil
}
