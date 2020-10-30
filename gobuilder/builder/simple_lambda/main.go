package main

import (
	"fmt"

	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	awslambda.Start(HandleEvent)
}

func HandleEvent(message Test) {
	fmt.Println(message.Message)

}

type Test struct {
	Message string `json:"message"`
}
