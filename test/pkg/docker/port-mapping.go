package docker

import (
	"fmt"
)

type PortMapping struct {
	host          string
	port          string
	containerPort string
	network       string
}

func NewPortMapping(host, port, containerPort, network string) (m PortMapping) {
	m = PortMapping{
		host:          host,
		port:          port,
		containerPort: containerPort,
		network:       network,
	}

	return
}

func (m PortMapping) Address() (address string) {
	const (
		addressFormat = "%s:%s"
	)

	address = fmt.Sprintf(addressFormat,
		m.host,
		m.port,
	)

	return
}

func (m PortMapping) Network() string {
	return m.network
}

func (m PortMapping) String() (s string) {
	const (
		portMappingFormat = "%s:%s:%s/%s"
	)

	s = fmt.Sprintf(portMappingFormat,
		m.host,
		m.port,
		m.containerPort,
		m.network,
	)

	return
}
