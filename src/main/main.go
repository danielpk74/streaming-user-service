package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"models"
	"services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)
	body := request.Body

	// Marshall the requrest body
	user := models.User{}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "error maping user object\n", StatusCode: http.StatusUnprocessableEntity}, err
	}

	_, err = services.CreateUser(user)
	if err != nil {
		fmt.Println("Got error calling create")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	// Log and return result
	fmt.Println("Wrote item:  ", &user)
	return events.APIGatewayProxyResponse{Body: "Success\n", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
