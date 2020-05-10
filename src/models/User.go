package models

import (
	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/danielpk74/service-core/security"
)

type User struct {
	ID       uint32 `json: "id"`
	Nickname string `json: "nickname" validate"required,gte=4,lte=20"`
	Email    string `json: "email" validate"email, required"`
	Password string `json: "password" validate"required"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
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
		if u.Nickname == "" {
			return errors.New("Required nickname")
		}
		if u.Email == "" {
			return errors.New("Required nickname")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
	default:
		if u.Nickname == "" {
			return errors.New("Required nickname")
		}
		if u.Password == "" {
			return errors.New("Required password")
		}
		if u.Email == "" {
			return errors.New("Required nickname")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
	}
	return nil
}
