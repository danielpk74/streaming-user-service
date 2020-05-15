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
	user := models.User{}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return Response{Body: "error maping user object\n", StatusCode: http.StatusUnprocessableEntity}, err
	}

	var out []byte
	var stCode int
	token, err := services.LoginUser(&user)
	if err != nil {
		fmt.Println("Got error getting the token", err.Error())
		out, _ = json.Marshal(&LoginFailResponse{
			Success: "false",
			Message: err.Error(),
		})
		stCode = 500
	} else {
		fmt.Println("Token: ", token)
		out, _ = json.Marshal(&LoginSuccessResponse{
			Success: "true",
			Token:   token,
		})
		stCode = 200
	}

	resp := Response{
		StatusCode:      stCode,
		IsBase64Encoded: false,
		Body:            string(out),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
