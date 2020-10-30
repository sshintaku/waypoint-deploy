package builder

import (
	"regexp"
	"testing"
)

func TestLambdaBuilder(t *testing.T) {
	var b Builder
	b.config.OutputName = "lambda"
	b.config.Source = "./"
	var lambda Lambda
	lambda.Amd64 = true
	lambda.Linux = true
	b.config.Lambda = &lambda
	c := BuildCommand(&b)
	checkCommand := "GOARCH=amd64 GOOS=linux go build -o " + b.config.OutputName + " " + b.config.Source
	output := c.String()
	matched, _ := regexp.MatchString(checkCommand, output)
	if matched == false {
		t.Fatal()
	}
}

func TestCommandBuilder(t *testing.T) {
	var b Builder
	b.config.OutputName = "main"
	b.config.Source = "./go-app"
	c := BuildCommand(&b)
	output := c.String()
	checkCommand := "go build -o " + b.config.OutputName + " " + b.config.Source
	matched, _ := regexp.MatchString(checkCommand, output)
	if matched == false {
		t.Fatal()
	}
}
