package services

import (
	"errors"
	"fmt"

	"models"
	"repository/crud"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUser(user *models.User) (*dynamodb.PutItemOutput, error) {
	user.Prepare()
	user.HashPassword()
	err := user.Validate("")
	if err != nil {
		return nil, errors.New("User object is not valid for creating.")
	}

	fmt.Println("SERVICE: User Prepared & validated: ", user)
	return crud.Save(user)
}

func LoginUser(user *models.User) (models.User, error) {
	user.Prepare()
	user.HashPassword()
	err := user.Validate("login")
	if err != nil {
		return models.User{}, errors.New("User object is not valid to login.")
	}
	return crud.LoginUser(user.Email, user.Password)
}
