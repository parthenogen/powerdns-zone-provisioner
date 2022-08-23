package docker

import (
	"context"
	"os"
	"os/exec"
)

type Container struct {
	command *exec.Cmd
	cancel  context.CancelFunc
	errors  chan error
}

func NewContainer(imageRef, name string) (c *Container, e error) {
	const (
		commandName = "docker"
		commandArg0 = "run"
		commandArg1 = "--name"
		commandArg3 = "--rm"
	)

	var (
		contextCancelFunc context.CancelFunc
		contextWithCancel context.Context
	)

	contextWithCancel, contextCancelFunc = context.WithCancel(
		context.Background(),
	)

	c = &Container{
		command: exec.CommandContext(contextWithCancel,
			commandName,
			commandArg0,
			commandArg1,
			name,
			commandArg2,
			imageRef,
		),
		cancel: contextCancelFunc,
		errors: make(chan error),
	}

	return
}

func (c *Container) Run() {
	go c.run()

	return
}

func (c *Container) run() {
	var (
		e error
	)

	e = c.command.Run()

	c.errors <- e

	return
}

func (c *Container) Error() error {
	return <-c.errors
}

func (c *Container) Stop() {
	c.cancel()

	return
}
