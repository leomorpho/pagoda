package tester

import (
	"context"
	"testing"

	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/enttest"
	"github.com/stretchr/testify/assert"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// https://testcontainers.com/guides/getting-started-with-testcontainers-for-go/
// https://golang.testcontainers.org/modules/postgres/
func CreateTestContainerPostgresConnStr(t *testing.T) (string, context.Context) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	return connStr, ctx
}

func CreateTestContainerPostgresEntClient(t *testing.T) (*ent.Client, context.Context) {
	connStr, ctx := CreateTestContainerPostgresConnStr(t)

	// Initialize Ent client with a test schema.
	client := enttest.Open(t, "postgres", connStr)
	t.Cleanup(func() {
		client.Close()
	})

	err := client.Schema.Create(ctx)
	assert.NoError(t, err)
	return client, ctx
}
