package platform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type RuntimeType int
type RegionType int

var RegionTypeArray []string = []string{}

var RuntimeArray []string = []string{
	"go1.x", "nodejs12.x", "nodejs10.x", "python3.8", "python3.7", "python3.6", "python2.7", "ruby2.7", "ruby2.5", "java11", "java8.al2", "java8", "dotnetcore3.1", "dotnetcore2.1", "provided.al2", "provided",
}

func ValidRuntimeValue() string {
	var result string
	for _, value := range RuntimeArray {
		result = result + value + "\n"
	}
	return result
}

func IsValidRuntime(element string) RuntimeType {

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

func (p Platform) CheckAccessPoint(accesspoint string) bool {
	client := efs.New(p.config.AwsSession)
	var input efs.DescribeAccessPointsInput
	input.SetAccessPointId(accesspoint)
	_, resultError := client.DescribeAccessPoints(&input)
	if resultError != nil {
		return false
	}
	return true
}

func (p Platform) KmsExists(kmsarn string) bool {
	client := kms.New(p.config.AwsSession)
	var input kms.DescribeKeyInput
	input.SetKeyId(kmsarn)
	_, resultError := client.DescribeKey(&input)
	if resultError != nil {
		return false
	}
	return true
}

func (p Platform) IamRoleExists(rolename string) bool {
	mysession, _ := CreateSessionWithRegion(p.config.Region)
	iamApi := iam.New(mysession)
	inputArray := strings.Split(rolename, "/")
	var input iam.GetRoleInput
	role := inputArray[len(inputArray)-1]
	input.RoleName = &role
	result, resultError := iamApi.GetRole(&input)
	if resultError != nil {
		return false
	}
	fmt.Println(result)
	return true
}

func (p Platform) IsValidSecurityGroup(secid string) bool {
	client := ec2.New(p.config.AwsSession)
	var input ec2.DescribeSecurityGroupReferencesInput
	input.GroupId = append(input.GroupId, &secid)
	result, secGroupError := client.DescribeSecurityGroupReferences(&input)
	if secGroupError != nil {
		return false
	}
	fmt.Println(result)
	return true
}

func (p Platform) IsValidSubnet(subnetid string) bool {
	client := ec2.New(p.config.AwsSession)
	var input ec2.DescribeSubnetsInput
	input.SubnetIds = append(input.SubnetIds, &subnetid)
	result, subnetError := client.DescribeSubnets(&input)
	if subnetError != nil {
		return false
	}
	fmt.Println(result)
	return true
}

func (p Platform) IsValidSnsTopic(topic_arn string) bool {
	client := sns.New(p.config.AwsSession)
	input := sns.GetTopicAttributesInput{
		TopicArn: &topic_arn,
	}
	result, resultError := client.GetTopicAttributes(&input)
	if resultError != nil {
		return false
	}
	fmt.Println(result)
	return true
}

func (p Platform) IsValidSqs(sqs_arn string) bool {
	client := sqs.New(p.config.AwsSession)
	var input sqs.ListQueuesInput
	arnSplit := strings.Split(sqs_arn, ":")
	queueName := arnSplit[len(arnSplit)-1]
	input.QueueNamePrefix = &queueName
	result, _ := client.ListQueues(&input)
	if result.QueueUrls == nil {
		return false
	}
	fmt.Println(result)
	return true
}

func (p Platform) LambdaFunctionExisits() bool {

	client := lambda.New(p.config.AwsSession)
	input := lambda.GetFunctionInput{
		FunctionName: aws.String(p.config.Lambda.FunctionName),
	}
	if p.config.Lambda.FunctionVersion != nil {
		input.Qualifier = aws.String(*p.config.Lambda.FunctionVersion)
	}
	result, resultError := client.GetFunction(&input)
	if resultError != nil {
		return false
	}
	fmt.Println(result)
	return true
}

func IsValidMountPath(mountpath string) bool {
	matched, _ := regexp.MatchString(`^/mnt/`, mountpath)
	return matched
}

func (p Platform) IsValidLambdaLayer(layer_arn string) bool {
	arnArray := strings.Split(layer_arn, ":")
	var arn string
	for i := 0; i < len(arnArray)-1; i++ {
		if i != 0 {
			arn = arn + ":" + arnArray[i]
		} else {
			arn = arn + arnArray[i]
		}
	}
	version := arnArray[len(arnArray)-1]
	arn_version, _ := strconv.ParseInt(version, 10, 64)
	mysession, _ := CreateSessionWithRegion(p.config.Region)
	client := lambda.New(mysession)
	var input lambda.GetLayerVersionInput
	input.LayerName = &arn
	input.VersionNumber = &arn_version
	version_output, layerError := client.GetLayerVersion(&input)
	if layerError != nil {
		return false
	}
	fmt.Println(version_output)
	return true

}
