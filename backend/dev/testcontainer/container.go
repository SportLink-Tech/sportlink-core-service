package testcontainer

import (
	"context"
	"time"

	"testing"

	"github.com/testcontainers/testcontainers-go"
)

func SportLinkContainer(t *testing.T, ctx context.Context) testcontainers.Container {
	// Create a context with timeout for container startup
	// The wait strategy has its own timeout, but we also need a context timeout
	containerCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	return createContainer(t, containerCtx, containerRequest())
}

func GetContainerEndpoint(t *testing.T, container testcontainers.Container, ctx context.Context) string {
	endpoint, err := container.PortEndpoint(ctx, "4566/tcp", "http")
	if err != nil {
		t.Fatalf("Failed to get container endpoint: %s", err)
	}
	return endpoint
}

func createContainer(t *testing.T, ctx context.Context, containerRequest testcontainers.ContainerRequest) testcontainers.Container {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start LocalStack: %s", err)
	}
	return container
}
