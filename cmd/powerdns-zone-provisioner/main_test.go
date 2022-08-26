package main

import (
	_ "embed"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/parthenogen/powerdns-zone-provisioner/test"
	"github.com/parthenogen/powerdns-zone-provisioner/test/pkg/dns"
)

var (
	//go:embed zones.yml
	testZoneFileYAML string
)

func TestMain(t *testing.T) {
	const (
		goModRoot = "../.."

		timeout = time.Second * 3

		zoneName = "example.com."

		rrSetName   = "www." + zoneName
		rrSetRecord = "192.0.2.1"

		tempFileDir = ""
		tempPattern = "*"
	)

	var (
		envVarKey string
		envVarVal string
		envVars   map[string]string

		auth   *test.PowerDNSAuth
		client *dns.Client

		zoneFile     *os.File
		zoneFileYAML string

		answer []string

		e error
	)

	auth, e = test.NewPowerDNSAuth(goModRoot, timeout)
	if e != nil {
		t.Fatal(e)
	}

	defer auth.Stop()

	zoneFile, e = os.CreateTemp(tempFileDir, tempPattern)
	if e != nil {
		t.Fatal(e)
	}

	defer os.Remove(
		zoneFile.Name(),
	)

	zoneFileYAML = fmt.Sprintf(testZoneFileYAML, zoneName, rrSetName, rrSetRecord)

	_, e = zoneFile.Write(
		[]byte(zoneFileYAML),
	)
	if e != nil {
		t.Fatal(e)
	}

	envVars = map[string]string{
		"SERVER_HOST": auth.IPAddress(),
		"SERVER_PORT": auth.PortAPI(),
		"API_KEY":     auth.APIKey(),
		"SERVER_ID":   auth.ServerID(),
		"ZONE_FILE":   zoneFile.Name(),
	}

	for envVarKey, envVarVal = range envVars {
		e = os.Setenv(envVarKey, envVarVal)
		if e != nil {
			t.Fatal(e)
		}
	}

	main()

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

	t.Log(
		auth.Stop(),
	)
}
