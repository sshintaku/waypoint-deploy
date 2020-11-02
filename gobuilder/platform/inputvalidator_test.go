package platform

import (
	"fmt"
	"testing"
)

func TestRuntimeInputTest(t *testing.T) {
	t.Run("Bad Input Test", func(t *testing.T) {
		indexCheck := IsValidRuntime("go")
		if indexCheck != -1 {
			t.Fatal("Error: return value was ")
		}
	})
	t.Run("Good Input Test", func(t *testing.T) {
		indexCheck2 := IsValidRuntime("go1.x")
		if indexCheck2 == -1 {
			t.Fatal("Error:  Could not find the string in the array")
		}
	})
}

func TestMemorySizeTest(t *testing.T) {
	t.Run("Low Memory Size Test", func(t *testing.T) {
		err1 := LambdaMemorySize(64)
		if err1 == nil {
			t.Fatal("Error: Memory size was in increments of 64 with a minimum of 128 and yet the test failed.")
		}
	})

	t.Run("Memory Divisible by 64", func(t *testing.T) {
		err1 := LambdaMemorySize(158)
		if err1 == nil {
			t.Fatal("Error: Memory size was not in increments of 64 with a minimum of 128 and yet it still passed.")
		}
	})
}

func TestQueueExists(t *testing.T) {
	var p Platform
	mysession, _ := CreateSessionWithRegion("us-east-1")
	p.config.AwsSession = mysession

	//var arn string
	arn := "arn:aws:sqs:us-east-1:003559363051:DeadLetterQueue"
	result := p.IsValidSqs(arn)
	if result == false {
		t.Fatal("Error: valid SQS queue reported as non-existant.")
	}
	//tarn := "SeadLetterQueue"

	fmt.Println(result)
}

func TestRoleExists(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(platform.config.Region)
	platform.config.AwsSession = client
	t.Run("Role exists confirmation", func(t *testing.T) {
		input := "arn:aws:iam::003559363051:role/service-role/SeijiTest-role-v83xhnlu"
		result := platform.IamRoleExists(input)
		if result == false {
			t.Fatal("IAM Role that should exist, didn't exist.")
		}
	})
	t.Run("Role not exist confirmation", func(t *testing.T) {
		input := "arn:aws:iam::003559363051:role/service-role/XeijiTest-role-v83xhnlu"
		result := platform.IamRoleExists(input)
		if result == true {
			t.Fatal("IAM Role that was expected to NOT exist does exist")
		}
	})
}

func TestLambdaLayerArn(t *testing.T) {
	var p Platform
	p.config.Region = "us-east-1"
	result := p.IsValidLambdaLayer("arn:aws:lambda:us-east-1:003559363051:layer:ToLambdaLayer")
	if result == false {
		t.Fatal("Error: Valid lambda layer arn was not found in the aws account.")
	}
}

func TestKMSExists(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	kmsarn := "arn:aws:kms:us-east-1:003559363051:key/c8962b7f-acd1-4ee8-a7e5-5ae19396be07"
	result := platform.KmsExists(kmsarn)
	if result != true {
		t.Fatal("Errro:  KMS key was not found when it should have been.")
	}
	fakekms := "arn:aws:kms:us-east-1:003559363051:key/c8962b7f-acd1-4ee8-a7e5-123456"
	result2 := platform.KmsExists(fakekms)
	if result2 == true {
		t.Fatal("Error:  Fake KMS key was check and it returned back as found when it should not have.")
	}
}

func TestSecurityGroupId(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	securityId := "sg-002a3033497278abd"
	result := platform.IsValidSecurityGroup(securityId)
	if result == false {
		t.Fatal("Error: Invalid Security ID.")
	}
	fakeSecurityId := "fake-security-id"
	result2 := platform.IsValidSecurityGroup(fakeSecurityId)
	if result2 == true {
		t.Fatal("Error:  Invalid Security ID checked and the result was the ID was valid.")
	}
}

func TestSubnetIs(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	subnetid := "subnet-fe1fc5b5"
	result := platform.IsValidSubnet(subnetid)
	if result == false {
		t.Fatal("Error: Invalid Network ID.")
	}
}

func TestSnsExists(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	topic := "arn:aws:sns:us-east-1:003559363051:SeijiTopic"
	result := platform.IsValidSnsTopic(topic)
	if result != true {
		t.Fatalf("Error: Valid topic showed up as false.")
	}
}

func TestSqsExists(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	sqs := "arn:aws:sqs:us-east-1:003559363051:DeadLetterQueue"
	result := platform.IsValidSqs(sqs)
	if result != true {
		t.Fatalf("Error: Valid SQS queue showed up as false.")
	}
	sqs2 := "arn:aws:sqs:us-east-1:003559363051:SeadLetterQueue"
	result2 := platform.IsValidSqs(sqs2)
	if result2 == true {
		t.Fatalf("Error: Invalid SQS queue showed up as valid.")
	}
}

func TestFunctionExists(t *testing.T) {
	var platform Platform
	var config DeployConfig
	config.Region = "us-east-1"
	platform.config = config
	function := "SeijiTest"
	client, _ := CreateSessionWithRegion(config.Region)
	platform.config.AwsSession = client
	platform.config.Lambda.FunctionName = function
	result := platform.LambdaFunctionExisits()
	if result != true {
		t.Fatalf("Error: Valid Lambda Function that should have existed was not found.")
	}
}

func TestLocalMountPath(t *testing.T) {
	var mnt = "/mnt/"
	if IsValidMountPath(mnt) == false {
		t.Fatal("Error: Invalid path.")
	}
	var mnt2 = "mnt/"
	if IsValidMountPath(mnt2) == true {
		t.Fatal("Error: Invalid path.")
	}
	var mnt3 = "//mnt"
	if IsValidMountPath(mnt3) == true {
		t.Fatal("Error: Invalid path.")
	}
	var mnt4 = "/us/mnt"
	if IsValidMountPath(mnt4) == true {
		t.Fatal("Error: Invalid path.")
	}
}
