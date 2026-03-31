package testcontainer

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func containerRequest() testcontainers.ContainerRequest {
	projectRoot, err := getProjectRoot()
	if err != nil {
		panic(err)
	}

	coreDynamoTablePath := filepath.Join(projectRoot, "backend", "dev", "localstack", "cloudformation", "core-dynamo-table.yml")
	sqsPath := filepath.Join(projectRoot, "backend", "dev", "localstack", "cloudformation", "sqs-queue.yml")
	initAwsPath := filepath.Join(projectRoot, "backend", "dev", "docker", "scripts")

	return testcontainers.ContainerRequest{
		Image:        "localstack/localstack:3.0.2",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES":       "sqs,dynamodb,ssm,cloudformation",
			"DEFAULT_REGION": "us-east-1",
			"AWS_REGION":     "us-east-1",
		},
		WaitingFor: wait.ForLog("All services created and database seeded.").
			WithStartupTimeout(2 * time.Minute),
		Mounts: []testcontainers.ContainerMount{
			testcontainers.BindMount(sqsPath, "/opt/code/localstack/sqs-queue.yml"),
			testcontainers.BindMount(coreDynamoTablePath, "/opt/code/localstack/core-dynamo-table.yml"),
			testcontainers.BindMount(initAwsPath, "/etc/localstack/init/ready.d"),
			testcontainers.BindMount("/var/run/docker.sock", "/var/run/docker.sock"),
		},
	}
}

func getProjectRoot() (string, error) {
	// Obtiene la ruta actual
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	// Navega hacia arriba hasta encontrar el directorio raíz del proyecto
	for {
		if strings.HasSuffix(currentDir, "sportlink-core-service") {
			return currentDir, nil
		}
		currentDir = filepath.Dir(currentDir)
	}
}
