package provisioners

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/mittwald/go-powerdns/apis/zones"
	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/docker"
)

func TestProvisioner(t *testing.T) {
	const (
		ipAddress = "127.0.0.1"
		portAPI   = "8081"
		portDNS   = "5353"

		addressAPI = ipAddress + ":" + portAPI
		addressDNS = ipAddress + ":" + portDNS

		apiKey           = "test"
		baseURL          = "http://" + addressAPI
		buildContextPath = "../.."
		dockerfilePath   = "../../test/build/powerdns-auth/Dockerfile"
		imageRef         = "test-provisioner"
		networkAPI       = "tcp"
		networkDNS       = "udp"
		portMappingAPI   = addressAPI + ":" + portAPI
		portMappingDNS   = addressDNS + ":" + "53/udp"
		serverID         = "localhost"
		timeout          = time.Second

		zoneName = "example.com."

		rrSetName   = "www.example.com."
		rrSetRecord = "192.0.2.1"
		rrSetTTL    = 60
		rrSetType   = "A"
	)

	var (
		dialer net.Dialer
		server *docker.Container

		p *provisioner

		ctx  context.Context
		zone zones.Zone

		client  dns.Client
		ipAddr  net.IPAddr
		message *dns.Msg

		e error
	)

	e = docker.Build(buildContextPath, dockerfilePath, imageRef)
	if e != nil {
		t.Fatal(e)
	}

	server, e = docker.NewContainer(imageRef,
		docker.WithPortMapping(portMappingAPI),
		docker.WithPortMapping(portMappingDNS),
	)
	if e != nil {
		t.Fatal(e)
	}

	server.Run()

	for {
		_, e = dialer.Dial(networkAPI, addressAPI)
		if e == nil {
			break
		}
	}

	p, e = NewProvisioner(baseURL, apiKey, serverID)
	if e != nil {
		t.Fatal(e)
	}

	zone = zones.Zone{
		Name: zoneName,
		ResourceRecordSets: []zones.ResourceRecordSet{
			{
				Name: rrSetName,
				Type: rrSetType,
				TTL:  rrSetTTL,
				Records: []zones.Record{
					{
						Content: rrSetRecord,
					},
				},
			},
		},
	}

	ctx, _ = context.WithTimeout(
		context.Background(),
		timeout,
	)

	e = p.Provision(zone, ctx)
	if e != nil {
		t.Fatal(e)
	}

	for {
		_, e = dialer.Dial(networkDNS, addressDNS)
		if e == nil {
			break
		}
	}

	message = new(dns.Msg)

	message.SetQuestion(rrSetName, dns.TypeA)

	message, _, e = client.Exchange(message, addressDNS)
	if e != nil {
		t.Fatal(e)
	}

	if len(message.Answer) != 1 {
		t.Fail()
	}

	ipAddr = net.IPAddr{
		IP: message.Answer[0].(*dns.A).A,
	}

	if ipAddr.String() != rrSetRecord {
		t.Fail()
	}

	server.Stop()

	e = server.Error()
	if e != nil {
		t.Log(e)
	}
}
