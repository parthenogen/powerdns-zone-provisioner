package docker

import (
	"os"
	"os/exec"
	"time"

	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/transport"
)

type Container struct {
	command *exec.Cmd
	errors  chan error

	portMappings []PortMapping
}

func NewContainer(imageRef string, options ...containerOption) (
	c *Container, e error,
) {
	const (
		commandName = "docker"
	)

	var (
		option containerOption
	)

	c = &Container{
		errors: make(chan error),
	}

	for _, option = range options {
		e = option(c)
		if e != nil {
			return
		}
	}

	c.command = exec.Command(commandName,
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
		commandArg2 = "--init"

		publishFlag = "-p"
	)

	var (
		portMapping PortMapping
	)

	commandArgs = []string{
		commandArg0,
		commandArg1,
		commandArg2,
	}

	for _, portMapping = range c.portMappings {
		commandArgs = append(commandArgs,
			publishFlag,
			portMapping.String(),
		)
	}

	commandArgs = append(commandArgs, imageRef)

	return
}

func (c *Container) Run() {
	go c.run()

	return
}

func (c *Container) RunAndDial(dialerTimeout time.Duration) (e error) {
	var (
		dialer      *transport.Dialer
		portMapping PortMapping
	)

	go c.run()

	dialer = transport.NewDialer()

	for _, portMapping = range c.portMappings {
		e = dialer.Dial(
			portMapping.Network(),
			portMapping.Address(),
			dialerTimeout,
		)
		if e != nil {
			return
		}
	}

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
	c.command.Process.Signal(os.Interrupt)

	return
}

type containerOption func(*Container) error

func WithPortMapping(host, port, containerPort, network string) (
	o containerOption,
) {
	o = func(c *Container) (e error) {
		c.portMappings = append(c.portMappings,
			NewPortMapping(host, port, containerPort, network),
		)

		return
	}

	return
}
