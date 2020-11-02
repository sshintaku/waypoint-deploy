package platform

// Comment out the following functions to implement Authentication in
// your component

//func (p *Platform) ValidateAuthFunc() interface{} {
//	return p.validateAuth
//}
//
//// AuthFunc satisfies the Authenticator interface
//func (p *Platform) AuthFunc() interface{} {
//	return r.authenticate
//}

// A ValidateAuthFunc does not have a strict signature, you can define the parameters
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
//
// If an error is returned, Waypoint will attempt to call
// AuthFunc
//func (p *Platform) validateAuth(
//	ctx context.Context,
//	ui terminal.UI,
//) error {
//	s := ui.Status()
//	defer s.Close()
//	s.Update("Validate authentication")
//
//	return fmt.Errorf("Unable to Authenticate")
//}

// A AuthFunc does not have a strict signature, you can define the parameters
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
//
// Output parameters must be *component.AuthResult, error
//func (p *Platform) authenticate(
//	ctx context.Context,
//	ui terminal.UI,
//) (*component.AuthResult, error) {
//
//	ui.Output("Describe the manual authentication steps here")
//
//	return &component.AuthResult{Authenticated: false}, nil
//}

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
func CreateSessionWithRegion(region string) (*session.Session, error) {
	client, clientError := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	return client, clientError
}
