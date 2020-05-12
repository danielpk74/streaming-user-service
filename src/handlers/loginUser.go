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

	// Marshall the request body
	user := models.User{}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "error maping user object\n", StatusCode: http.StatusUnprocessableEntity}, err
	}

	token, err := services.LoginUser(&user)
	if err != nil {
		fmt.Println("Got error getting the token")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: http.StatusUnprocessableEntity}, nil
	}

	// Log and return result
	fmt.Println("Token success", token)
	return events.APIGatewayProxyResponse{Body: "Success\n", StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
