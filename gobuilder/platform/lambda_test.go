package platform

import (
	"testing"
)

func TestLambdaBuilder(t *testing.T) {
	var platform Platform
	var lambda LambdaInput
	lambda.FunctionName = "Lambda GOLANG"
	lambda.Description = "GoLambdaTest"
	//lambda.HandlerName = "lambda"
	platform.config.SourceBinary = "./lambda"
	platform.config.Lambda.RoleArn = "arn:aws:iam::003559363051:role/service-role/SeijiTest-role-v83xhnlu"

	platform.config.Lambda = lambda
	zipError := ZipCreationFunction(&platform)
	if zipError != nil {
		t.Fatal(zipError)
	}
	lambdaError := CreateLambda(&platform)
	if lambdaError != nil {

	}
	//fmt.Println(output)
}
