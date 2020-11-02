package platform

import (
	"testing"
)

func TestZipDirectory(t *testing.T) {
	var platform Platform
	var config DeployConfig
	var directory string
	directory = "./python"
	config.LambdaFiles.SourceFolder = &directory
	platform.config = config
	platform.ZipDirectoryFiles()
}
