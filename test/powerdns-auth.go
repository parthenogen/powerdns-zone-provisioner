package test

import (
	"path"
	"time"

	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/docker"
)

const (
	apiKey    = "test"
	ipAddress = "127.0.0.1"
	portAPI   = "8081"
	portDNS   = "5353"
	serverID  = "localhost"
)

type PowerDNSAuth struct {
	container *docker.Container
}

func NewPowerDNSAuth(goModRoot string, dialTimeout time.Duration) (
	a *PowerDNSAuth, e error,
) {
	const (
		containerPortAPI = "8081"
		containerPortDNS = "53"
		dockerfilePath   = "test/build/powerdns-auth/Dockerfile"
		imageRef         = "test-powerdns-auth"
		networkAPI       = "tcp"
		networkDNS       = "udp"
	)

	e = docker.Build(
		goModRoot,
		path.Join(goModRoot, dockerfilePath),
		imageRef,
	)
	if e != nil {
		return
	}

	a = &PowerDNSAuth{}

	a.container, e = docker.NewContainer(imageRef,
		docker.WithPortMapping(ipAddress, portAPI, containerPortAPI, networkAPI),
		docker.WithPortMapping(ipAddress, portDNS, containerPortDNS, networkDNS),
	)
	if e != nil {
		return
	}

	e = a.container.RunAndDial(dialTimeout)
	if e != nil {
		return
	}

	return
}

func (a *PowerDNSAuth) APIKey() string {
	return apiKey
}

func (a *PowerDNSAuth) ServerID() string {
	return serverID
}

func (a *PowerDNSAuth) IPAddress() string {
	return ipAddress
}

func (a *PowerDNSAuth) PortAPI() string {
	return portAPI
}

func (a *PowerDNSAuth) PortDNS() string {
	return portDNS
}

func (a *PowerDNSAuth) Stop() (e error) {
	a.container.Stop()

	return a.container.Error()
}
