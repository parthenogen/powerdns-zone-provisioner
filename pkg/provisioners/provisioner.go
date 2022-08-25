package provisioners

import (
	"context"

	"github.com/mittwald/go-powerdns"
	"github.com/mittwald/go-powerdns/apis/zones"
)

type provisioner struct {
	client   zones.Client
	serverID string
}

func NewProvisioner(baseURL, apiKey, serverID string) (
	p *provisioner, e error,
) {
	var (
		client pdns.Client
	)

	client, e = pdns.New(
		pdns.WithBaseURL(baseURL),
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

func (p *provisioner) Provision(zone zones.Zone, ctx context.Context) (
	e error,
) {
	_, e = p.client.CreateZone(
		ctx,
		p.serverID,
		zone,
	)
	if e != nil {
		return
	}

	return
}
