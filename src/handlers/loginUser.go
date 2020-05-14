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

type Response events.APIGatewayProxyResponse

type LoginSuccessResponse struct {
	Success string
	Token   string
}

type LoginFailResponse struct {
	Success string
	Message string
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("Received body: ", request.Body)
	body := request.Body

	// Marshall the request body
	user := models.User{}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return Response{Body: "error maping user object\n", StatusCode: http.StatusUnprocessableEntity}, err
	}

	var b []byte
	token, err := services.LoginUser(&user)
	if err != nil {
		fmt.Println("Got error getting the token", err.Error())
		b, _ = json.Marshal(&LoginFailResponse{
			Success: "false",
			Message: err.Error(),
		})
	} else {
		fmt.Println("Token: ", token)
		b, _ = json.Marshal(&LoginSuccessResponse{
			Success: "true",
			Token:   token,
		})
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(b),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
