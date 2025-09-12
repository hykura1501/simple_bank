package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return errors.New("must contain only lowercase letters, digits, underscore")
	}
	return nil
}

func ValidateFullName(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(username) {
		return errors.New("must contain only letters, spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 5, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("is not a valid email address")
	}
	return nil
}
