package platform

import (
	"fmt"
	"testing"

	awslambda "github.com/aws/aws-sdk-go/service/lambda"
)

func TestLambdaBuilder(t *testing.T) {
	var platform Platform
	var lambda LambdaInput
	var config DeployConfig
	config.Region = "us-east-1"

	lambda.FunctionName = "Lambda GOLANG"
	lambda.Description = "GoLambdaTest"
	//lambda.HandlerName = "lambda"
	//config. = "./lambda"
	var role string
	role = "SeijiTest-role-v83xhnlu"
	platform.config.Lambda.RoleArn = &role

	platform.config.Lambda = lambda
	platform.config = config
	zipError := platform.ZipCreationFunction()
	if zipError != nil {
		fmt.Println("")
	}
	//platform.config.Lambda.ZipFileInBytes = &zipBytes

	lambdaError := platform.CreateLambda()
	if lambdaError != nil {

	}
	//fmt.Println(output)
}

func TestLambdaUpdate(t *testing.T) {
	var platform Platform
	platform.config.Region = "us-east-1"
	platform.config.Lambda.FunctionName = "SeijiFunctionTest"
	var revision string
	revision = "b256a131-3606-47e6-8190-4bba358ef3be"
	platform.config.Lambda.RevisionId = &revision
	var dryrun bool
	dryrun = false
	platform.config.Lambda.DryRun = &dryrun
	var file string
	file = "./lambda"
	platform.config.LambdaFiles.SourceBinary = &file
	platform.config.Lambda.Publish = true
	client, _ := CreateSessionWithRegion(platform.config.Region)
	platform.config.AwsLambda = awslambda.New(client)
	platform.ZipCreationFunction()
	result := platform.UpdateLambdaFunction()
	fmt.Println(result)
}

func TestLambdaCreation(t *testing.T) {
	var p Platform
	p.config.Region = "us-east-1"
	mysession, sessionError := CreateSessionWithRegion(p.config.Region)
	if sessionError != nil {

	}
	p.config.AwsSession = mysession
	var source string
	source = "./lambda"

	var arn string
	arn = "arn:aws:iam::003559363051:role/service-role/SeijiTest-role-v83xhnlu"
	var kms string
	kms = "arn:aws:kms:us-east-1:003559363051:key/c8962b7f-acd1-4ee8-a7e5-5ae19396be07"
	p.config.Lambda.KMSKeyArn = &kms
	p.config.LambdaFiles.SourceBinary = &source
	p.config.Lambda.FunctionName = "NewGoLambda"
	p.config.Lambda.RoleArn = &arn
	p.config.Lambda.Runtime = "go1.x"
	p.config.Lambda.Publish = false
	var dl DeadLetterConfig
	var accessarn string
	accessarn = "arn:aws:sqs:us-east-1:003559363051:DeadLetterQueue"
	dl.TargetArn = &accessarn
	p.config.Lambda.DeadLetter = &dl
	var fs FileSystemConfig
	var fsystem string
	//fsystem = "arn:aws:elasticfilesystem:us-east-1:003559363051:access-point/fsap-0cabb6c2bf376680f"
	fsystem = "arn:aws:elasticfilesystem:us-east-1:003559363051:access-point/fsap-0cabb6c2bf376680f"
	fs.Arn = &fsystem
	var mpath string
	mpath = "/mnt/efs0"
	fs.LocalMountPath = &mpath

	p.config.Lambda.FileSystemConfigs = append(p.config.Lambda.FileSystemConfigs, &fs)
	zipError := p.ZipCreationFunction()
	if zipError != nil {
		fmt.Println("")
	}
	var memSize int64
	memSize = 192
	p.config.Lambda.MemorySize = &memSize
	var vpc Vpc
	vpc.SecurityIds = append(vpc.SecurityIds, "sg-03a975d66845c08e0")
	vpc.Subnets = append(vpc.Subnets, "subnet-494cf62d")
	p.config.Lambda.VpcConfig = &vpc
	env := make(map[string]*string)
	var value string
	value = "Seiji"
	env["owner"] = &value
	p.config.Lambda.Environment = env
	tag := make(map[string]*string)
	var tagValue string
	tagValue = "Project X"
	tag["Project"] = &tagValue
	p.config.Lambda.Tags = tag
	layerArn := "arn:aws:lambda:us-east-1:003559363051:layer:ToLambdaLayer:1"
	p.config.Lambda.Layers = append(p.config.Lambda.Layers, &layerArn)
	//p.config.Lambda.ZipFileInBytes = &zipBytes
	lambdaError := p.CreateLambda()
	if lambdaError != nil {
		fmt.Println(lambdaError)
		t.Fatal(lambdaError)
	}

}
