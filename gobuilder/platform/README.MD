# Lambda WayPoint Plugin

Waypoint allows flexibility and extendability by means of plugins that are written in the GOLANG programming language.  All plugins use the HCL format to parse and injest paramters and interfaces to be consumed by the plugin to understand what the user wants to create and where.  The HCL file requires multiple inputs.  The plugin works very dilligently to confirm those items exist in the cloud environments before it tries to create / update a Lambda function.

## Expected Inputs

The HCL input file should follow the following structure for the plugin to injest:

```
type DeployConfig struct {
	Region       string           `hcl:"region"`
	SourceBinary string           `hcl:"binary"`
	Lambda       LambdaInput      `hcl:"lambda"`
	AwsSession   *session.Session `hcl:"awssession"`
}

type LambdaInput struct {
	Description       string              `hcl:"description,optional"`
	FunctionName      string              `hcl:"functioname"`
	RoleArn           *string             `hcl:"rolearn"`
	DeadLetter        *DeadLetterConfig   `hcl:"deadletterqueue,optional"`
	Environment       map[string]*string  `hcl:"environment,optional"`
	FileSystemConfigs []*FileSystemConfig `hcl:"filesystemconfig,optional"`
	KMSKeyArn         *string             `hcl:"kmskeyarn,optional"`
	Layers            []*string           `hcl:"layers,optional"`
	MemorySize        *int64              `hcl:"memorysize,optional"`
	Publish           *bool               `hcl:"publish,optional"`
	Runtime           string              `hcl:"runtime"`
	Tags              map[string]*string  `hcl:"tags,optional"`

	/* Default timeout value is 3 seconds.  Maximum allowed is 900.  This number is in metrics of seconds
	Timeout       *int64         `hcl:"timeout,optional"` */
	TracingConfig *TracingConfig `hcl:"tracingconfig,optional"`
	VpcConfig     *Vpc           `hcl:"vpcconfig,optional"`
}

type Vpc struct {
	Subnets     []string `hcl:"subnets"`
	SecurityIds []string `hcl:"securityids"`
}

type TracingConfig struct {
	Mode string `hcl:"mode"`
}

type FileSystemConfig struct {
	FilesystemName string `hcl:"filesystem_name"`
	LocalMountPath string `hcl:"localmountpath"`
}

type DeadLetterConfig struct {
	TargetArn *string `hcl:"targetarn, optional"`
}
```

### AWS Authentication
Your AWS access tokens can be set in multiple ways: 1) Environment variables, 2) Shared configuration files.  Ensure that your tokens are set either way and the plugin will gather the AWS tokens from either of those two methods.  Ensure that in your HCL file, the property of Region is populated.  The plugin needs that to determine which datacenter to targe the lambda function.  From the tokens and the region, the application has enough information to create a session with your AWS subscription.


### Dead Letter Configuration

A Dead Letter Configuration is either an SNS topic or SQS queue.  Even though the name of the properyt is "TargetArn", it really is **not** expecting an ARN.  In fact it needs to one of two things: SNS Topic or the SQS URL.  If the "TargetArn" property is set, the plugin will try to determine if the property is a fully formed URL.  If it is, it will query your AWS account for that SQS queue.  If it does not exist, the application will throw and error stating that the input that was specified is invalid.  

If the TargetArn is **not** a URL, it will treat the value as an SNS Topic.  It will query your subscription to ensure that the topic exists.  If the topic does not exist, it will throw and error stating that the input that was specified is invalid.  

### EFS Filesystems

If a list of EFS Filesystems are specified in the HCL, the plugin will first confirm that a.) the file system does exist in your subscription / region and b.) that a filesystem mount point is specified.  If either of those conditions are not true, the plugin will throw an error stating that either the filesystem does **not** exist, the filesystem exists but the mount point wasn't specified, or that the filesystem does **not** exist *AND* a mount point was not specified.

### AWS Role for the Lambda function uses

The Lambda needs a role to use to authorize the function to run.  This role needs to be specified in the HCL file.  The plugin will verify that the role actually does exist before it tries to create the lambda function.  If role doesn't exist, it will throw an error stating that the specified role was not found in your subscription / region specified.

### AWS Functions

The plugin does check to see if the Function exists in the region you are specifying.  This is **case sensitive** so "TestFunction" is different from "testfunction".  There is also a property on the object that is a boolean value.  This must be set to True for the Lambda code to be updated.  If it is not set, the plugin will throw an error and will stop immediately. Setting the property "Publish" will increment the version push out the code immediately.