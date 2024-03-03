package data

import (
	"fmt"
	"regexp"

	"github.com/EZCampusDevs/firepit/util"
)

const (
	PASSWORD_MIN_LENGTH = 4
	PASSWORD_MAX_LENGTH = 72

	USERNAME_MIN_LENGTH = 4
	USERNAME_MAX_LENGTH = 16
)

var (
	UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")
)

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthPayload) IsValidUsername() error {

	if util.IsEmptyOrWhitespace(a.Username) {

		return fmt.Errorf("Username cannot be empty")
	}

	if len(a.Username) < USERNAME_MIN_LENGTH || len(a.Username) > USERNAME_MAX_LENGTH {

		return fmt.Errorf("Username must have length between %d and %d", USERNAME_MIN_LENGTH, USERNAME_MAX_LENGTH)
	}

	if !UsernameRegex.MatchString(a.Username) {

		return fmt.Errorf("Username contains invalid characters")
	}

	return nil
}

func (a *AuthPayload) IsValidPassword() error {

	if util.IsEmptyOrWhitespace(a.Password) {

		return fmt.Errorf("Password cannot be empty")
	}

	if len(a.Password) < PASSWORD_MIN_LENGTH || len(a.Password) > PASSWORD_MAX_LENGTH {

		return fmt.Errorf("Password must be between %d and %d", PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH)
	}

	return nil
}

func (a *AuthPayload) IsValid() error {

	if err := a.IsValidUsername(); err != nil {
		return err
	}

	if err := a.IsValidPassword(); err != nil {
		return err
	}

	return nil
}
