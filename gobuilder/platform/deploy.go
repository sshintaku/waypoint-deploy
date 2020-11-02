package platform

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/waypoint-plugin-examples/gobuilder/registry"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type DeployConfig struct {
	Region      string           `hcl:"region"`
	Lambda      LambdaInput      `hcl:"lambda"`
	AwsSession  *session.Session `hcl:"awssession"`
	AwsLambda   *lambda.Lambda   `hcl:"awslambda,optional"`
	LambdaFiles FileInformation  `hcl:"filelocation"`
}

type FileInformation struct {
	SourceBinary *string    `hcl:"sourcebinary,optional"`
	SourceFolder *string    `hcl:"sourcefolder,optional"`
	Layer        *LayerInfo `hcl:"layer,optional"`
}

type LayerInfo struct {
	FolderLocation string `hcl:"folderlocation"`
	Runtime        string `hcl:"runtimetype"`
}

type LambdaInput struct {
	Description       string              `hcl:"description,optional"`
	DryRun            *bool               `hcl:"dryrun,optional"`
	FunctionName      string              `hcl:"functioname"`
	FunctionVersion   *string             `hcl:"functionversion,optional"`
	RoleArn           *string             `hcl:"rolearn,optional"`
	RevisionId        *string             `hcl:"revisionid,optional"`
	DeadLetter        *DeadLetterConfig   `hcl:"deadletterqueue,optional"`
	Environment       map[string]*string  `hcl:"environment,optional"`
	FileSystemConfigs []*FileSystemConfig `hcl:"filesystemconfig,optional"`
	KMSKeyArn         *string             `hcl:"kmskeyarn,optional"`
	Layers            []*string           `hcl:"layers,optional"`
	MemorySize        *int64              `hcl:"memorysize,optional"`
	Publish           bool                `hcl:"publish"`
	Runtime           string              `hcl:"runtime"`
	Tags              map[string]*string  `hcl:"tags,optional"`
	UpdateLambda      *bool               `hcl:"updatelambda,optional"`

	/* Default timeout value is 3 seconds.  Maximum allowed is 900.  This number is in metrics of seconds
	Timeout       *int64         `hcl:"timeout,optional"` */
	TracingConfig  *TracingConfig `hcl:"tracingconfig,optional"`
	VpcConfig      *Vpc           `hcl:"vpcconfig,optional"`
	ZipFileInBytes *[]byte        `hcl:"zipfileinbytes,optional"`
}

type Vpc struct {
	Subnets     []string `hcl:"subnets"`
	SecurityIds []string `hcl:"securityids"`
}

type TracingConfig struct {
	Mode string `hcl:"mode"`
}

type FileSystemConfig struct {
	Arn            *string `hcl:"arn"`
	LocalMountPath *string `hcl:"localmountpath"`
}

type DeadLetterConfig struct {
	TargetArn *string `hcl:"targetarn, optional"`
}

type Platform struct {
	config DeployConfig
}

// Implement Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// Implement ConfigurableNotify
func (p *Platform) ConfigSet(config interface{}) error {
	// c, ok := config.(*DeployConfig)
	// if !ok {
	// 	// The Waypoint SDK should ensure this never gets hit
	// 	return fmt.Errorf("Expected *DeployConfig as parameter")
	// }

	// validate the config
	// if c.Region == "" {
	// 	return fmt.Errorf("Region must be set to a valid directory")
	// }

	return nil
}

// Implement Builder
func (p *Platform) DeployFunc() interface{} {
	// return a function which will be called by Waypoint
	return p.deploy
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet

// In addition to default input parameters the registry.Artifact from the Build step
// can also be injected.
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (p *Platform) deploy(ctx context.Context, ui terminal.UI, artifact *registry.Artifact, src *component.Source, job *component.JobInfo) (*Deployment, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Validating Lambda inputs")
	// Validation Step
	u.Step(terminal.StatusOK, "Lambda input validation is complete.  Creating zip file of the application.")
	zipError := p.ZipCreationFunction()
	if zipError != nil {
		u.Step(terminal.ErrorBoldStyle, "Zipfile creation experienced a fata error.")
	}
	//p.config.Lambda.ZipFileInBytes = &zipBytes

	//utils.DefaultSubnets(ctx, sess)

	return &Deployment{}, nil
}
