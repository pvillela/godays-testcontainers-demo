package util

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

// LaunchPostgres launches a Postgres container and returns a handle to it and its connection URL.
func LaunchPostgres(ctx context.Context) (postgres tc.Container, postgresUrl string) {
	log.Println("Starting postgres container...")
	postgresPort := nat.Port("5432/tcp")
	postgres, err := tc.GenericContainer(ctx,
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgres",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "pass",
					"POSTGRES_USER":     "user",
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true, // auto-start the container
		})
	if err != nil {
		log.Fatal("start:", err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatal("map:", err)
	}
	postgresURLTemplate := "postgres://user:pass@localhost:%s?sslmode=disable"
	postgresUrl = fmt.Sprintf(postgresURLTemplate, hostPort.Port())
	log.Printf("Postgres container started, running at:  %s\n", postgresUrl)
	return
}
