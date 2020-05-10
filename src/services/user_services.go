package services

import (
	"errors"

	"models"
	"repository/crud"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUser(user *models.User) (*dynamodb.PutItemOutput, error) {
	user.Prepare()
	err := user.Validate("")
	if err != nil {
		return nil, errors.New("Post not found")
	}

	return crud.Save(user)
}
