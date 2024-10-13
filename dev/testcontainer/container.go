package testcontainer

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

func SportLinkContainer(t *testing.T, ctx context.Context) testcontainers.Container {
	return createContainer(t, ctx, containerRequest())
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
