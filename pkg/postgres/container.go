package postgres

import (
	"context"
	"fmt"

	"git.h2hsecure.com/pda/commons/pkg/logs"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgres testcontainers.Container
)

func CreatePostgresqlDb(ctx context.Context) (*string, *int, error) {
	log.Info().
		Str("ctx", "postgres").
		Msg("Postgresql database starting...")

	postgresPort := nat.Port("5432/tcp")

	req := testcontainers.ContainerRequest{
		Image: "registry.h2hsecure.com/pda/postgresqldb:latest",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
		},
		ExposedPorts: []string{string(postgresPort)},
		WaitingFor:   wait.ForAll(wait.ForExposedPort()),
		Hostname:     "postgresqldb",
	}

	var err error

	postgres, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     testcontainers.ProviderDefault,
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("postgres container start %w", err)
	}

	if err := postgres.StartLogProducer(ctx); err != nil {
		return nil, nil, fmt.Errorf("postgres log producer: %w", err)
	}

	postgres.FollowOutput(logs.NewContainerLogger("postgres"))

	postgresHost, err := postgres.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres container host: %w", err)
	}

	postgresPorts, err := postgres.MappedPort(ctx, postgresPort)

	if err != nil {
		return nil, nil, fmt.Errorf("postgres container host: %w", err)
	}
	port := postgresPorts.Int()

	return &postgresHost, &port, nil
}

func PostgresTerminate(ctx context.Context) error {
	if err := postgres.StopLogProducer(); err != nil {
		return fmt.Errorf("stop logging: %w", err)
	}

	if err := postgres.Terminate(ctx); err != nil {
		return fmt.Errorf("stop logging: %w", err)
	}

	return nil
}
