package platform

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
)

func (p Platform) CreateLambda() error {

	client := awslambda.New(p.config.AwsSession)
	p.config.AwsLambda = client
	info, fileStat := os.Stat(*p.config.LambdaFiles.SourceBinary)
	if fileStat != nil {
		return fileStat
	}
	//  This section of checking Roles does not work.
	//  I need to focus on this section
	/* roleExists := p.IamRoleExists(*p.config.Lambda.RoleArn)
	if roleExists == false {
		return fmt.Errorf("Error: IAM Role Arn does not exist.")
	} */
	runtimeExists := IsValidRuntime(p.config.Lambda.Runtime)
	if runtimeExists == -1 {
		message := ValidRuntimeValue()
		return fmt.Errorf("Invalid Runtime specified.  Valid values are the following: \n\n%s", message)
	}

	if p.config.Lambda.MemorySize == nil {
		var size int64
		size = 128
		p.config.Lambda.MemorySize = &size
	}
	zipLocation := *p.config.LambdaFiles.SourceBinary + ".zip"
	zipBytes, _ := ioutil.ReadFile(zipLocation)
	//  The following are the bare minimum properties
	lambdaInput := awslambda.CreateFunctionInput{
		Code: &awslambda.FunctionCode{
			ZipFile: zipBytes,
		},
		Description:  aws.String(p.config.Lambda.Description),
		FunctionName: aws.String(p.config.Lambda.FunctionName),
		Role:         aws.String(*p.config.Lambda.RoleArn),
		Handler:      aws.String(info.Name()),
		Runtime:      aws.String(p.config.Lambda.Runtime),
		Publish:      aws.Bool(p.config.Lambda.Publish),
	}
	if p.config.Lambda.DeadLetter != nil {
		isSqs := p.IsValidSqs(*p.config.Lambda.DeadLetter.TargetArn)
		isSns := p.IsValidSnsTopic(*p.config.Lambda.DeadLetter.TargetArn)
		if isSns == true {
			snsResult := p.IsValidSnsTopic(*p.config.Lambda.DeadLetter.TargetArn)
			if snsResult == false {
				return fmt.Errorf("Error: Dead Letter Queue is neither a valid SQS Url or SNS Topic.")
			} else {
				lambdaInput.DeadLetterConfig.TargetArn = p.config.Lambda.DeadLetter.TargetArn
			}
		}
		if isSqs == true {
			var dl awslambda.DeadLetterConfig
			var target string
			target = *p.config.Lambda.DeadLetter.TargetArn
			dl.SetTargetArn(target)
			lambdaInput.DeadLetterConfig = &dl
		}
		if isSns == false && isSqs == false {
			return fmt.Errorf("Error: SQS or SNS Arn was not a valid ARN.  ARN checked was the following: %s", *p.config.Lambda.DeadLetter.TargetArn)
		}
	}
	if p.config.Lambda.VpcConfig != nil {

		var input awslambda.VpcConfig
		for _, secgroup := range p.config.Lambda.VpcConfig.SecurityIds {
			result := p.IsValidSecurityGroup(secgroup)
			if result == false {
				return fmt.Errorf("Error: Invalid Security Group ID specified.  Security Group ID is: %v", secgroup)
			}
			input.SecurityGroupIds = append(input.SecurityGroupIds, &secgroup)
		}
		for _, subnet := range p.config.Lambda.VpcConfig.Subnets {
			result := p.IsValidSubnet(subnet)
			if result == false {
				return fmt.Errorf("Error: Invalid Subnet ID specified.  Subnet ID is: %v", subnet)
			}
			input.SubnetIds = append(input.SubnetIds, &subnet)
		}
		lambdaInput.VpcConfig = &input
	}

	if p.config.Lambda.FileSystemConfigs != nil {
		if p.config.Lambda.VpcConfig == nil {
			return fmt.Errorf("Error:  VPC subnet and security ids must be defined to mount EFS filesystems.")
		} else {
			for _, fs := range p.config.Lambda.FileSystemConfigs {
				var fsObject awslambda.FileSystemConfig
				if fs.LocalMountPath == nil {
					return fmt.Errorf("Error:  You cannot have a filesystem input and not have a mountpath input.")
				} else {
					var fsarn *string
					fsarn = fs.Arn
					fsObject.Arn = fsarn
					var lmp *string
					lmp = fs.LocalMountPath
					fsObject.LocalMountPath = lmp
					if IsValidMountPath(*lmp) == false {
						return fmt.Errorf("Error:  Local Mount Path is invalid.  All mount points must be valid and start with /mnt/ as the path.")
					}
					var fs awslambda.FileSystemConfig
					fs.SetArn(*fsarn)
					fs.SetLocalMountPath(*lmp)
					lambdaInput.FileSystemConfigs = append(lambdaInput.FileSystemConfigs, &fs)
				}
			}
		}
	}

	if p.config.Lambda.TracingConfig != nil {
		lambdaInput.TracingConfig.Mode = aws.String(p.config.Lambda.TracingConfig.Mode)
	}

	if p.config.Lambda.Layers != nil {
		var list []*string
		for _, v := range p.config.Lambda.Layers {
			result := p.IsValidLambdaLayer(*v)
			if result == true {
				list = append(list, v)
				lambdaInput.Layers = list
			}
			if result == false {
				return fmt.Errorf("Error:  The targe Lambda Layer ARN is either invalid or poorly formed.  It needs both the lambda arn needs to have the version at the end of the arn.  Lambda Layer ARN is: %v", *v)
			}
		}
	}
	if p.config.Lambda.Environment != nil {
		var env awslambda.Environment
		env.Variables = p.config.Lambda.Environment
		lambdaInput.Environment = &env
	}

	if p.config.Lambda.MemorySize != nil {
		memoryCheck := LambdaMemorySize(*p.config.Lambda.MemorySize)
		if memoryCheck != nil {
			return fmt.Errorf("Error: Lambda memory size needs to be at least 128 and in increments of 64.")
		}
		lambdaInput.MemorySize = p.config.Lambda.MemorySize

	}

	if p.config.Lambda.Tags != nil {
		//tags := p.CreateTags()
		lambdaInput.Tags = p.config.Lambda.Tags
	}

	if p.config.Lambda.KMSKeyArn != nil {
		result := p.KmsExists(*p.config.Lambda.KMSKeyArn)
		if result == true {
			lambdaInput.KMSKeyArn = aws.String(*p.config.Lambda.KMSKeyArn)
		} else {
			return fmt.Errorf("Error: Lambda KMS input isn't a valid ARN.")
		}
	}

	_, errLambda := p.config.AwsLambda.CreateFunction(&lambdaInput)
	if errLambda != nil {
		return errLambda
	}
	return nil

}

/* func (p Platform) CreateLambdaInput() error {
	mySession, sessionError := CreateSessionWithRegion(p.config.Region)
	if sessionError != nil {
		return fmt.Errorf("Error: There was a problem creating an AWS session in the Create Lambda function %v", sessionError)
	}

	// Create a Lambda client from just a session.
	client := awslambda.New(mySession)
	p.config.AwsLambda = client

	runtimeCheck := IsValidRuntime(p.config.Lambda.Runtime)
	if runtimeCheck == -1 {
		return fmt.Errorf("Error: The input value for the runtime is invalid.  Please check the runtime is a valid Lambda runtime version.")
	}
	return nil
} */

func (p Platform) UpdateLambdaFunction() error {
	info, fileStat := os.Stat(*p.config.LambdaFiles.SourceBinary)
	if fileStat != nil {
		return fileStat
	}

	zipBytes, readBinaryError := ioutil.ReadFile(info.Name() + ".zip")
	if readBinaryError != nil {
		return readBinaryError
	}

	input := awslambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(p.config.Lambda.FunctionName),
		ZipFile:      zipBytes,
		Publish:      aws.Bool(p.config.Lambda.Publish),
	}
	if p.config.Lambda.RevisionId != nil {
		input.RevisionId = aws.String(*p.config.Lambda.RevisionId)
	}

	if p.config.Lambda.DryRun != nil {
		input.DryRun = aws.Bool(*p.config.Lambda.DryRun)
	}
	validate := input.Validate()
	fmt.Println(validate)
	result, resultError := p.config.AwsLambda.UpdateFunctionCode(&input)
	if resultError != nil {
		return resultError
	}
	fmt.Println(result)
	return nil
}

/* func CreateTestFunction(p Platform) error {
	session := CreateSession()
	client := awslambda.New(session)
	p.config.AwsLambda = client
	/* info, fileStat := os.Stat(*p.config.LambdaFiles.SourceBinary)
	if fileStat != nil {
		return fileStat
	}
	zipError := p.ZipCreationFunction()

	if zipError != nil {
		return fmt.Errorf("Zip Error ")
	}
	zipBytes, _ := ioutil.ReadFile("./lambda.zip")

	lambdaInput := &awslambda.CreateFunctionInput{
		Code: &awslambda.FunctionCode{
			ZipFile: zipBytes,
		},
		Description:  aws.String("Test Function"),
		FunctionName: aws.String("GoTestFunction"),
		Role:         aws.String("arn:aws:iam::003559363051:role/service-role/SeijiTest-role-v83xhnlu"),
		Handler:      aws.String("lambda"),
		Runtime:      aws.String("go1.x"),
	}
	_, errLambda := client.CreateFunction(lambdaInput)
	if errLambda != nil {
		return errLambda
	}
	return nil
} */
