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

	portMappings []string // e.g., {"127.0.0.1:5353:53/udp"}
}

func NewContainer(imageRef string, options ...containerOption) (
	c *Container, e error,
) {
	const (
		commandName = "docker"
	)

	var (
		contextCancelFunc context.CancelFunc
		contextWithCancel context.Context

		option containerOption
	)

	contextWithCancel, contextCancelFunc = context.WithCancel(
		context.Background(),
	)

	c = &Container{
		cancel: contextCancelFunc,
		errors: make(chan error),
	}

	for _, option = range options {
		e = option(c)
		if e != nil {
			return
		}
	}

	c.command = exec.CommandContext(contextWithCancel,
		commandName,
		c.makeCommandArgs(imageRef)...,
	)

	c.command.Stdout = os.Stdout
	c.command.Stderr = os.Stderr

	return
}

func (c *Container) makeCommandArgs(imageRef string) (commandArgs []string) {
	const (
		commandArg0 = "run"
		commandArg1 = "--rm"

		publishFlag = "-p"
	)

	var (
		portMapping string
	)

	commandArgs = []string{
		commandArg0,
		commandArg1,
	}

	for _, portMapping = range c.portMappings {
		commandArgs = append(commandArgs, publishFlag, portMapping)
	}

	commandArgs = append(commandArgs, imageRef)

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

type containerOption func(*Container) error

func WithPortMapping(portMapping string) (o containerOption) {
	o = func(c *Container) (e error) {
		c.portMappings = append(c.portMappings, portMapping)

		return
	}

	return
}
