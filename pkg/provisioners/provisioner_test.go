package provisioners

import (
	"testing"
	"time"

	"github.com/mittwald/go-powerdns/apis/zones"
	"github.com/parthenogen/powerdns-zone-provisioner/test"
	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/dns"
)

func TestProvisioner(t *testing.T) {
	const (
		goModRoot = "../.."

		timeout = time.Second * 3

		zoneName = "example.com."

		rrSetName   = "www." + zoneName
		rrSetRecord = "192.0.2.1"
		rrSetTTL    = 60
		rrSetType   = "A"
	)

	var (
		auth *test.PowerDNSAuth

		p    *provisioner
		zone zones.Zone

		answer []string
		client *dns.Client

		e error
	)

	auth, e = test.NewPowerDNSAuth(goModRoot, timeout)
	if e != nil {
		t.Fatal(e)
	}

	defer auth.Stop()

	p, e = NewProvisioner(
		auth.IPAddress(),
		auth.PortAPI(),
		auth.APIKey(),
		auth.ServerID(),
		timeout,
	)
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

	e = p.Provision(zone)
	if e != nil {
		t.Fatal(e)
	}

	client = dns.NewClient()

	answer, e = client.QueryTypeA(
		auth.IPAddress(),
		auth.PortDNS(),
		rrSetName,
	)
	if e != nil {
		t.Fatal(e)
	}

	if len(answer) != 1 {
		t.Fail()
	}

	if answer[0] != rrSetRecord {
		t.Fail()
	}

	// test re-provisioning of existing zone

	e = p.Provision(zone)
	if e != nil {
		t.Fatal(e)
	}

	answer, e = client.QueryTypeA(
		auth.IPAddress(),
		auth.PortDNS(),
		rrSetName,
	)
	if e != nil {
		t.Fatal(e)
	}

	if len(answer) != 1 {
		t.Fail()
	}

	if answer[0] != rrSetRecord {
		t.Fail()
	}

	e = auth.Stop()
	if e != nil {
		t.Log(e)
	}
}
