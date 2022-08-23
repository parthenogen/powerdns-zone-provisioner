package docker

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/parthenogen/redis-cluster/test/pkg/docker"
)

func TestContainer(t *testing.T) {
	const (
		ipAddress = "127.0.0.1"
		port      = "8081"

		address          = ipAddress + ":" + port
		apiKey           = "test"
		buildContextPath = "../../.."
		dockerfilePath   = "../../build/powerdns-auth/Dockerfile"
		imageRef         = "test-container"
		network          = "tcp"
		portMapping      = address + ":" + port
		url              = "http://" + address + "/api"
		xAPIKeyHeaderKey = "X-API-Key"

		expectedResponse = `[{"url": "/api/v1", "version": 1}]`
	)

	var (
		container *Container

		client         http.Client
		dialer         net.Dialer
		request        *http.Request
		response       *http.Response
		responseBytes  []byte
		responseString string

		e error
	)

	e = docker.Build(buildContextPath, dockerfilePath, imageRef)
	if e != nil {
		t.Fatal(e)
	}

	container, e = NewContainer(imageRef,
		WithPortMapping(portMapping),
	)
	if e != nil {
		t.Fatal(e)
	}

	container.Run()

	for {
		_, e = dialer.Dial(network, address)
		if e == nil {
			break
		}
	}

	request, e = http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		t.Fatal(e)
	}

	request.Header.Add(xAPIKeyHeaderKey, apiKey)

	response, e = client.Do(request)
	if e != nil {
		t.Fatal(e)
	}

	if response.StatusCode != http.StatusOK {
		t.Fail()
	}

	defer response.Body.Close()

	responseBytes, e = io.ReadAll(response.Body)
	if e != nil {
		t.Fatal(e)
	}

	responseString = strings.TrimSpace(
		string(responseBytes),
	)

	if responseString != expectedResponse {
		t.Fail()
	}

	container.Stop()

	e = container.Error()
	if e != nil {
		t.Log(e)
	}
}
