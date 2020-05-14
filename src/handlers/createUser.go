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

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("Received body: ", request.Body)

	body := request.Body
	user := models.User{}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		return Response{Body: "error maping user object\n", StatusCode: http.StatusUnprocessableEntity}, err
	}

	_, err = services.CreateUser(&user)
	if err != nil {
		fmt.Println("Got error calling create")
		fmt.Println(err.Error())
		return Response{Body: "Error", StatusCode: http.StatusUnprocessableEntity}, nil
	}

	user.Password = ""
	out, _ := json.Marshal(user)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(out),
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
