package platform

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
)

func CreateLambda(p *Platform) error {

	mySession := CreateSession()

	// Create a Lambda client from just a session.
	svc := awslambda.New(mySession)

	info, fileStat := os.Stat(p.config.SourceBinary)
	if fileStat != nil {
		return fileStat
	}
	lambdaBytes, readBinaryError := ioutil.ReadFile(info.Name() + ".zip")
	if readBinaryError != nil {

		return readBinaryError
	}
	runtimeCheck := RuntimeIndex(p.config.Lambda.Runtime)
	if runtimeCheck == -1 {
		return fmt.Errorf("Error: The input value for the runtime is invalid.  Please check the runtime is a valid Lambda runtime version.")
	}

	roleExists := IamRoleExists(p.config.Lambda.RoleArn)
	if roleExists == false {
		return fmt.Errorf("Error: IAM Role Arn does not exist.")
	}

	lambdaInput := &awslambda.CreateFunctionInput{
		Code: &awslambda.FunctionCode{
			ZipFile: lambdaBytes,
		},
		Description:  &p.config.Lambda.Description,
		FunctionName: aws.String(p.config.Lambda.FunctionName),
		Role:         aws.String(p.config.Lambda.RoleArn),
		Handler:      aws.String(info.Name()),
		Runtime:      aws.String(p.config.Lambda.Runtime),
	}
	if p.config.Lambda.Publish != nil {
		lambdaInput.Publish = p.config.Lambda.Publish
	}
	if p.config.Lambda.DeadLetter != nil && p.config.Lambda.DeadLetter.TargetArn != nil {
		result := IsValidArn(*p.config.Lambda.DeadLetter.TargetArn)
		if result != true {
			return fmt.Errorf("Error:  The targe DeadLetter ARN is invalid.  Target ARN is: %v", p.config.Lambda.DeadLetter.TargetArn)
		}
	}

	if p.config.Lambda.FileSystemConfigs != nil {
		for _, fs := range p.config.Lambda.FileSystemConfigs {
			result := IsValidEfsVolue(fs.Arn)
			if result == false {
				return fmt.Errorf("Error: Invalid EFS Volume ARN specified.  EFS ARN is: %v", fs.Arn)
			}
			var fsObject lambda.FileSystemConfig
			fsObject.Arn = &fs.Arn
			fsObject.LocalMountPath = &fs.LocalMountPath
			lambdaInput.FileSystemConfigs = append(lambdaInput.FileSystemConfigs, &fsObject)
		}
	}

	if p.config.Lambda.Layers != nil {
		var list []string
		for _, v := range p.config.Lambda.Layers {
			result := IsValidArn(*v)
			if result == true {
				list = append(list, *v)
			}
			if result == false {
				return fmt.Errorf("Error:  The targe Lambda Layer ARN is invalid.  Lambda Layer ARN is: %v", *v)
			}
		}
	}
	if p.config.Lambda.Environment != nil {
		m := GetEnvironmentVariables(p)
		lambdaInput.Environment.Variables = aws.StringMap(*m)
	}

	if p.config.Lambda.MemorySize != nil {
		memoryCheck := LambdaMemorySize(*p.config.Lambda.MemorySize)
		if memoryCheck != nil {
			return fmt.Errorf("Error: Lambda memory size needs to be at least 128 and in increments of 64.")
		}
		lambdaInput.MemorySize = p.config.Lambda.MemorySize

	}

	if p.config.Lambda.Tags != nil {
		tags := CreateTags(p)
		lambdaInput.Tags = aws.StringMap(*tags)
	}

	if p.config.Lambda.KMSKeyArn != nil {
		result := doesKMSexist(*p.config.Lambda.KMSKeyArn)
		if result == true {
			lambdaInput.KMSKeyArn = aws.String(*p.config.Lambda.KMSKeyArn)
		} else {
			return fmt.Errorf("Error: Lambda KMS input isn't a valid ARN.")
		}
	}

	_, errLambda := svc.CreateFunction(lambdaInput)
	if errLambda != nil {
		return errLambda
	}
	return nil

}
