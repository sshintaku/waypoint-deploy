module github.com/hashicorp/waypoint-plugin-examples/gobuilder

go 1.14

require (
	github.com/aws/aws-sdk-go v1.35.18
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/waypoint v0.1.4
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20201021094150-1b1044b1478e
	github.com/mitchellh/go-glint v0.0.0-20201015034436-f80573c636de
	google.golang.org/protobuf v1.25.0
)

// replace github.com/hashicorp/waypoint-plugin-sdk => ../../waypoint-plugin-sdk
