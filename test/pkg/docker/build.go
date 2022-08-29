package docker

import (
	"os"
	"os/exec"
)

func Build(buildContextPath, dockerfilePath, imageRef string) (e error) {
	const (
		commandName = "docker"
		commandArg0 = "build"
		commandArg2 = "-f"
		commandArg4 = "-t"
	)

	var (
		command *exec.Cmd
	)

	command = exec.Command(commandName,
		commandArg0,
		buildContextPath,
		commandArg2,
		dockerfilePath,
		commandArg4,
		imageRef,
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	e = command.Run()
	if e != nil {
		return
	}

	return
}
