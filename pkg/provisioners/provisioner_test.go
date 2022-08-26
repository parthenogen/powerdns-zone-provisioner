package provisioners

import (
	"testing"
	"time"

	"github.com/mittwald/go-powerdns/apis/zones"
	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/dns"
	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/docker"
)

func TestProvisioner(t *testing.T) {
	const (
		ipAddress = "127.0.0.1"
		portAPI   = "8081"
		portDNS   = "5353"

		apiKey           = "test"
		buildContextPath = "../.."
		containerPortDNS = "53"
		dockerfilePath   = "../../test/build/powerdns-auth/Dockerfile"
		imageRef         = "test-provisioner"
		networkAPI       = "tcp"
		networkDNS       = "udp"
		serverID         = "localhost"
		timeout          = time.Second * 3

		zoneName = "example.com."

		rrSetName   = "www." + zoneName
		rrSetRecord = "192.0.2.1"
		rrSetTTL    = 60
		rrSetType   = "A"
	)

	var (
		server *docker.Container

		p    *provisioner
		zone zones.Zone

		answer []string
		client *dns.Client

		e error
	)

	e = docker.Build(buildContextPath, dockerfilePath, imageRef)
	if e != nil {
		t.Fatal(e)
	}

	server, e = docker.NewContainer(imageRef,
		docker.WithPortMapping(ipAddress, portAPI, portAPI, networkAPI),
		docker.WithPortMapping(ipAddress, portDNS, containerPortDNS, networkDNS),
	)
	if e != nil {
		t.Fatal(e)
	}

	e = server.RunAndDial(timeout)
	if e != nil {
		t.Fatal(e)
	}

	defer server.Stop()

	p, e = NewProvisioner(ipAddress, portAPI, apiKey, serverID)
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

	e = p.Provision(zone, timeout)
	if e != nil {
		t.Fatal(e)
	}

	client = dns.NewClient()

	answer, e = client.QueryTypeA(ipAddress, portDNS, rrSetName)
	if e != nil {
		t.Fatal(e)
	}

	if len(answer) != 1 {
		t.Fail()
	}

	if answer[0] != rrSetRecord {
		t.Fail()
	}

	server.Stop()

	e = server.Error()
	if e != nil {
		t.Log(e)
	}
}
