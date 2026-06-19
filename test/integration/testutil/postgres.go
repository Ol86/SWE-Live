//go:build integration

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
	DSN       string
}

func StartPostgres(t *testing.T, ctx context.Context) *PostgresContainer {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "library_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get postgres host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("get postgres port: %v", err)
	}

	return &PostgresContainer{
		Container: container,
		DSN:       fmt.Sprintf("postgres://test:test@%s:%s/library_test?sslmode=disable", host, port.Port()),
	}
}

func OpenPool(t *testing.T, ctx context.Context, dsn string) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open postgres pool: %v", err)
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}

	return pool
}

func ApplySchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	root := findModuleRoot(t)
	schema, err := os.ReadFile(filepath.Join(root, "internal", "db", "schema", "schema.sql"))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	if _, err := pool.Exec(ctx, string(schema)); err != nil {
		t.Fatalf("apply schema: %v", err)
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
