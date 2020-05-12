package crud

import (
	"fmt"
	"os"

	"models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/danielpk74/service-core/utils/channels"
)

func Save(user *models.User) (*dynamodb.PutItemOutput, error) {
	fmt.Println("REPOSITORY: User received: ", user)

	// Create the dynamo client object
	done := make(chan bool)
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	var err error
	var it *dynamodb.PutItemOutput

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("REPOSITORY: User prepared for Dynamo: ", av)
	go func(ch chan<- bool) {
		defer close(ch)
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		}

		it, err = svc.PutItem(input)

		if err != nil {
			ch <- false
			return
		}

		ch <- true

	}(done)

	if channels.OK(done) {
		return it, nil
	}

	return nil, err
}

func LoginUser(email, password string) (models.User, error) {
	fmt.Println("REPOSITORY: User & Pass: ", email, password)

	// Create the dynamo client object
	done := make(chan bool)
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	userdb := models.User{}
	var result *dynamodb.GetItemOutput
	var err error

	go func(ch chan<- bool) {
		defer close(ch)

		fmt.Println("Trying to read from table: ", os.Getenv("TABLE_NAME"))
		result, err = svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String(email),
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Unmarshall the result in to an Item
		err = dynamodbattribute.UnmarshalMap(result.Item, &userdb)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if err != nil {
			ch <- false
			return
		}

		ch <- true

	}(done)

	if channels.OK(done) {
		return userdb, nil
	}

	return models.User{}, err
}
