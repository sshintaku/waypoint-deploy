package platform

import (
	"testing"
)

func RuntimeInputTest(t *testing.T) {
	t.Run("Bad Input Test", func(t *testing.T) {
		indexCheck := RuntimeIndex("go")
		if indexCheck != -1 {
			t.Fatal("Error: return value was ")
		}
	})
	t.Run("Good Input Test", func(t *testing.T) {
		indexCheck2 := RuntimeIndex("go1.x")
		if indexCheck2 == -1 {
			t.Fatal("Error:  Could not find the string in the array")
		}
	})
}

func MemorySizeTest(t *testing.T) {
	t.Run("Low Memory Size Test", func(t *testing.T) {
		err1 := LambdaMemorySize(64)
		if err1 == nil {
			t.Fatal("Error: Unexpected Lambda Memory Size validation.")
		}
	})

	t.Run("Memory Divisible by 64", func(t *testing.T) {
		err1 := LambdaMemorySize(158)
		if err1 != nil {
			t.Fatal("Error: Memory should be in increments of 64 with a minimum of 128.")
		}
	})
}

func CheckRoleExists(t *testing.T) {
	t.Run("Role exists confirmation", func(t *testing.T) {
		var input *string
		*input = "arn:aws:iam::003559363051:role/service-role/SeijiTest-role-v83xhnlu"
		result := IamRoleExists(*input)
		if result == false {
			t.Fatal("IAM Role that was expected to exist does not exist")
		}
	})
	t.Run("Role not exist confirmation", func(t *testing.T) {
		var input *string
		*input = "arn:aws:iam::003559363051:role/service-role/SeijiTest-role"
		result := IamRoleExists(*input)
		if result == true {
			t.Fatal("IAM Role that was expected to NOT exist does exist")
		}
	})
}

func TestKMSExists(t *testing.T) {
	kmsarn := "arn:aws:kms:us-east-1:003559363051:key/c8962b7f-acd1-4ee8-a7e5-5ae19396be07"
	result := doesKMSexist(kmsarn)
	if result != true {
		t.Fatal("Errro:  KMS key was not found when it should have been.")
	}
	fakekms := "arn:aws:kms:us-east-1:003559363051:key/c8962b7f-acd1-4ee8-a7e5-123456"
	result2 := doesKMSexist(fakekms)
	if result2 == true {
		t.Fatal("Error:  Fake KMS key was check and it returned back as found when it should not have.")
	}
}

func TestArn(t *testing.T) {
	validArn := IsValidArn("arn:aws:sqs:us-east-1:003559363051:DeadLetterQueue")
	if validArn != true {
		t.Fatal("Error: SQS or SNS Arn provided is not valid")
	}
}
