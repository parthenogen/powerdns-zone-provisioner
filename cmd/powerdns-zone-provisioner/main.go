package main

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/mittwald/go-powerdns/apis/zones"
	"github.com/parthenogen/powerdns-zone-provisioner/pkg/provisioners"
	"gopkg.in/yaml.v3"
)

func main() {
	type environment struct {
		ServerHost string        `env:"SERVER_HOST"`
		ServerPort string        `env:"SERVER_PORT"`
		ApiKey     string        `env:"API_KEY"`
		ServerID   string        `env:"SERVER_ID" envDefault:"localhost"`
		Timeout    time.Duration `env:"TIMEOUT" envDefault:"3s"`
		ZoneFile   string        `env:"ZONE_FILE"`
	}

	type provisioner interface {
		Provision(zones.Zone) error
	}

	var (
		envVars environment
		p       provisioner

		zone     zones.Zone
		zoneFile []byte
		zoneList []zones.Zone

		e error
	)

	e = env.Parse(&envVars)
	if e != nil {
		log.Fatalln(e)
	}

	p, e = provisioners.NewProvisioner(
		envVars.ServerHost,
		envVars.ServerPort,
		envVars.ApiKey,
		envVars.ServerID,
		envVars.Timeout,
	)
	if e != nil {
		log.Fatalln(e)
	}

	zoneFile, e = os.ReadFile(envVars.ZoneFile)
	if e != nil {
		log.Fatalln(e)
	}

	e = yaml.Unmarshal(zoneFile, &zoneList)
	if e != nil {
		log.Fatalln(e)
	}

	for _, zone = range zoneList {
		e = p.Provision(zone)
		if e != nil {
			log.Fatalln(e)
		}
	}
}
