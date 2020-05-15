package services

import (
	"errors"
	"fmt"

	"models"
	"repository/crud"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	auth "github.com/danielpk74/service-core/auth"
	"github.com/google/uuid"
)

// Check if an email is already used by another account.
func validateIfEmailExists(e string) bool {
	var u *models.User
	u, _ = crud.GetUserByEmail(e)
	if u.Email != "" {
		return true
	}

	return false
}

func CreateUser(user *models.User) (*dynamodb.PutItemOutput, error) {
	user.Prepare()

	u := validateIfEmailExists(user.Email)
	if u {
		return nil, errors.New("The Email is already used.")
	}

	user.HashPassword()
	err := user.Validate("")
	if err != nil {
		return nil, errors.New("User object is not valid for creating.")
	}

	fmt.Println("SERVICE: User Prepared & validated: ", &user)
	return crud.Save(user)
}

func LoginUser(user *models.User) (string, error) {
	user.Prepare()
	var err error
	err = user.Validate("login")
	if err != nil {
		return "", errors.New("User object is not valid to login.")
	}

	var u *models.User
	u, err = crud.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}

	fmt.Println("User in DB", u)
	if u.Email == "" {
		return "", errors.New("User not found")
	}

	fmt.Println("Password from DB", u.Password)
	if user.ValidatePassword(u.Password) {
		fmt.Println("Password validated")
		return CreateToken(u.Id)
	}

	fmt.Println("Error after password validation")
	return "", err
}

func CreateToken(userId uuid.UUID) (string, error) {
	return auth.CreateToken(userId)
}
