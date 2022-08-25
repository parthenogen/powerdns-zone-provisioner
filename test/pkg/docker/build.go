package docker

import (
	"github.com/parthenogen/redis-cluster/test/pkg/docker"
)

var (
	Build func(string, string, string) error = docker.Build
)
