package pass

import (
	"errors"
	"fmt"
	"unicode"
)

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, encode string) (bool, error)
	RehashNeeded(encoded string) bool
}

func VerifyPassword(password string) error {
	var (
		uppercasePresent   = false
		lowercasePresent   = false
		numberPresent      = false
		specialCharPresent = false
		minPassLength      = 8
		maxPassLength      = 64
		passLen            int
	)
	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	if !lowercasePresent {
		return errors.New("lowercase letter missing")
	}
	if !uppercasePresent {
		return errors.New("uppercase letter missing")
	}
	if !numberPresent {
		return errors.New("at least one numeric character required")
	}
	if !specialCharPresent {
		return errors.New("special character missing")
	}
	if !(minPassLength <= passLen && passLen >= maxPassLength) {
		return fmt.Errorf("password length must be between %d to %d characters long", minPassLength, maxPassLength)
	}
	return nil
}
