package models

import (
	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/danielpk74/service-core/security"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username" validate:"required,gte=4,lte=20"`
	Email    string    `json:"email" validate:"email, required"`
	Password string    `json:"password,omitempty" validate:"required"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Id = uuid.New()
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
	case "update":
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
	default:
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Password == "" {
			return errors.New("Required password")
		}
		if u.Email == "" {
			return errors.New("Required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
	}
	return nil
}

func (u *User) ValidatePassword(passwordStored string) bool {
	err := security.VerifyPassword(passwordStored, u.Password)
	if err != nil {
		return false
	}

	return true
}
