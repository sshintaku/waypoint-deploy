package builder

import (
	"os"
	"os/exec"
)

func BuildCommand(b *Builder) *exec.Cmd {
	if b.config.Arch != nil && b.config.Arch.Amd64 == true {
		os.Unsetenv("GOARCH")
		os.Setenv("GOARCH", "amd64")
	}
	if b.config.Arch != nil && b.config.Arch.Linux == true {
		os.Unsetenv("GOOS")
		os.Setenv("GOOS", "linux")
	}
	return exec.Command(
		"go",
		"build",
		"-o",
		b.config.OutputName,
		b.config.Source,
	)
}
