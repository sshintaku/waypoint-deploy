package platform

import (
	"fmt"

	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
)

type RuntimeType int

var RuntimeArray []string = []string{
	"go1.x", "nodejs12.x", "nodejs10.x", "python3.8", "python3.7", "python3.6", "python2.7", "ruby2.7", "ruby2.5", "java11", "java8.al2", "java8", "dotnetcore3.1", "dotnetcore2.1", "provided.al2", "provided",
}

func RuntimeIndex(element string) RuntimeType {

	for k, v := range RuntimeArray {
		if element == v {
			var runtimeIndex RuntimeType = RuntimeType(k)
			return runtimeIndex
		}
	}
	return -1 //not found.
}

func LambdaMemorySize(memorysize int64) error {
	if memorysize < 128 {
		return fmt.Errorf("Error: The minimum value for lambda memory size is 128.")
	}
	divisible := memorysize % 64
	if divisible != 0 {
		return fmt.Errorf("Error: The minimum value for lambda memory size need to be in multiples of 64.")
	}
	return nil
}

func CreateTags(p *Platform) *map[string]string {
	var m map[string]string
	if p.config.Lambda.Tags == nil {
		return nil
	}
	for k, v := range p.config.Lambda.Tags {
		m[k] = *v
	}
	return &m
}

func GetEnvironmentVariables(p *Platform) *map[string]string {
	var m map[string]string
	if p.config.Lambda.Environment == nil {
		return nil
	}
	for k, v := range p.config.Lambda.Environment {
		m[k] = *v
	}
	return &m
}

func doesKMSexist(kmsarn string) bool {
	mysession := CreateSession()
	kmsSession := kms.New(mysession)
	var input kms.DescribeKeyInput
	input.SetKeyId(kmsarn)
	_, resultError := kmsSession.DescribeKey(&input)
	if resultError != nil {
		return false
	}
	return true
}

func IamRoleExists(rolename string) bool {
	mySession := CreateSession()
	iamApi := iam.New(mySession)
	var input iam.GetRoleInput
	input.RoleName = &rolename
	_, resultError := iamApi.GetRole(&input)
	if resultError != nil {
		return false
	}
	return true
}

func IsValidArn(arn string) bool {
	return awsarn.IsARN(arn)
}

func IsValidEfsVolue(arn string) bool {
	mysession := CreateSession()
	client := efs.New(mysession)
	var fsInput efs.DescribeFileSystemsInput
	fsInput.FileSystemId = &arn
	_, resultError := client.DescribeFileSystems(&fsInput)
	if resultError != nil {
		return false
	}
	return true
}
