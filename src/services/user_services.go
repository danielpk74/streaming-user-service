package services

import (
	"errors"
	"fmt"

	"models"
	"repository/crud"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	auth "github.com/danielpk74/service-core/auth"
)

func CreateUser(user *models.User) (*dynamodb.PutItemOutput, error) {
	user.Prepare()
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

	var dUser *models.User
	dUser, err = crud.LoginUser(user.Email, user.Password)
	if err != nil {
		return "", err
	}

	fmt.Println("User in DB", dUser)
	if dUser.Email == "" {
		return "", errors.New("User not found")
	}

	fmt.Println("Password from DB", dUser.Password)
	if user.ValidatePassword(dUser.Password) {
		fmt.Println("Password validated")
		return auth.CreateToken(dUser.Id)
	}

	fmt.Println("Error after password validation")
	return "", err
}
