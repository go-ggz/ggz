package integrations

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGGZServer(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "goggz/ggz-server",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog("Starting shorten server on :8080"),
	}
	ggzServer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatal(err)
	}

	// At the end of the test remove the container
	defer ggzServer.Terminate(ctx)
	// Retrieve the container IP
	ip, err := ggzServer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the port mapped to port 8080
	port, err := ggzServer.MappedPort(ctx, "8080")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/", ip, port.Port()))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}
