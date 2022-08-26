package provisioners

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mittwald/go-powerdns"
	"github.com/mittwald/go-powerdns/apis/zones"
)

type provisioner struct {
	client   zones.Client
	serverID string
}

func NewProvisioner(serverHost, serverPort, apiKey, serverID string) (
	p *provisioner, e error,
) {
	const (
		hostFormat = "%s:%s"
		scheme     = "http"
	)

	var (
		baseURL url.URL
		client  pdns.Client
	)

	baseURL = url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf(hostFormat, serverHost, serverPort),
	}

	client, e = pdns.New(
		pdns.WithBaseURL(
			baseURL.String(),
		),
		pdns.WithAPIKeyAuthentication(apiKey),
	)
	if e != nil {
		return
	}

	p = &provisioner{
		client:   client.Zones(),
		serverID: serverID,
	}

	return
}

func (p *provisioner) Provision(zone zones.Zone, timeout time.Duration) (
	e error,
) {
	var (
		timer context.Context
	)

	timer, _ = context.WithTimeout(
		context.Background(),
		timeout,
	)

	_, e = p.client.CreateZone(
		timer,
		p.serverID,
		zone,
	)
	if e != nil {
		return
	}

	return
}
