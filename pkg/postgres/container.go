package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresContainer *postgres.PostgresContainer
)

func CreatePostgresqlDb(ctx context.Context, net *testcontainers.DockerNetwork) (string, error) {
	log.Info().
		Str("ctx", "postgres").
		Msg("Postgresql database starting...")

	var err error
	postgresContainer, err = postgres.Run(ctx,
		"docker.io/postgres:17-alpine",
		network.WithNetwork([]string{"postgresql"}, net),
		postgres.WithDatabase("pda"),
		postgres.WithUsername("pda"),
		postgres.WithPassword("pda"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)

	if err != nil {
		return "", fmt.Errorf("postgres container start %w", err)
	}

	return postgresContainer.ConnectionString(ctx, "sslmode=disable")
}

func PostgresTerminate(ctx context.Context) error {
	if err := postgresContainer.StopLogProducer(); err != nil {
		return fmt.Errorf("stop logging: %w", err)
	}

	if err := postgresContainer.Terminate(ctx); err != nil {
		return fmt.Errorf("stop logging: %w", err)
	}

	return nil
}
