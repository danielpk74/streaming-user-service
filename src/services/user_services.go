package services

import (
	"errors"
	"fmt"

	"models"
	"repository/crud"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUser(user models.User) (*dynamodb.PutItemOutput, error) {
	user.Prepare()
	user.HashPassword()
	err := user.Validate("")
	if err != nil {
		return nil, errors.New("Post not found")
	}

	fmt.Println("SERVICE: User Prepared & validated: ", user)
	return crud.Save(user)
}
