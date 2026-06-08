package psql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

const (
	containerStartupTimeout = 60 * time.Second
	pingCheckInterval       = 500 * time.Millisecond
)

// postgresTmpfsCustomizer provides a ContainerCustomizer that mounts a tmpfs for
// /var/lib/postgresql/data, bypassing overlayfs permission restrictions that prevent
// the default postgres Docker image from initializing its data directory.
type postgresTmpfsCustomizer struct{}

func (p postgresTmpfsCustomizer) Customize(req *testcontainers.GenericContainerRequest) error {
	req.Tmpfs = map[string]string{"/var/lib/postgresql/data": "rw,noexec,nosuid,size=256m"}
	return nil
}

// waitForPostgresReady polls the container until PostgreSQL accepts connections via Ping.
func waitForPostgresReady(ctx context.Context, connStr string) error {
	ticker := time.NewTicker(pingCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for postgres: %w", ctx.Err())
		case <-ticker.C:
			db, err := sql.Open(string(Postgres), connStr)
			if err != nil {
				continue
			}
			if pingErr := db.PingContext(ctx); pingErr == nil {
				_ = db.Close()
				return nil
			}
			_ = db.Close()
		}
	}
}

// newPostgresContainer starts a PostgreSQL test container using the standard Docker entrypoint
// (which respects POSTGRES_HOST_AUTH_METHOD env vars) mounted on tmpfs, then returns the
// connection string and cleanup function.
func newPostgresContainer(ctx context.Context, t *testing.T) (string, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, containerStartupTimeout)
	defer cancel()

	container, err := testcontainers.Run(ctx, "postgres:16-alpine",
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":                 "testuser",
			"POSTGRES_PASSWORD":             "testpass",
			"POSTGRES_DB":                   "testdb",
			"POSTGRES_HOST_AUTH_METHOD":     "trust",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithCmd(
			"-c", "fsync=off",
			"-c", "synchronous_commit=off",
		),
		postgresTmpfsCustomizer{},
	)
	require.NoError(t, err, "failed to start postgres container")

	endpoint, err := container.PortEndpoint(ctx, "5432/tcp", "")
	require.NoError(t, err)
	connStr := fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", endpoint) // PortEndpoint already returns host:port

	waitErr := waitForPostgresReady(ctx, connStr)
	require.NoError(t, waitErr, "postgres did not become ready in time")

	return connStr, func() {
		_ = container.Terminate(ctx)
	}
}


// --- NewConnection success cases ---

func TestNewConnection_Success(t *testing.T) {
	ctx := context.Background()

	connStr, terminate := newPostgresContainer(ctx, t)
	defer terminate()

	client, err := NewConnection(ctx, connStr)
	require.NoError(t, err, "NewConnection should not return an error")

	assert.NotNil(t, client, "client should not be nil")
	assert.Equal(t, Postgres, Driver(client.driverName), "driver name should be pgx")
	assert.NotEmpty(t, client.connectionURI, "connection URI should not be empty")

	db := client.GetClient()
	assert.NotNil(t, db, "db field should not be nil")

	err = db.PingContext(ctx)
	assert.NoError(t, err, "PingContext on the underlying sqlx.DB should succeed")

	// Explicitly close before container teardown.
	err = client.Close()
	assert.NoError(t, err, "client.Close() should succeed")
}

func TestNewConnection_InvalidConnectionString(t *testing.T) {
	ctx := context.Background()

	client, err := NewConnection(ctx, "invalid-connection-string://nope")
	assert.Error(t, err, "NewConnection should return an error for malformed connection string")
	assert.Nil(t, client, "client should be nil when connection fails")
}

func TestNewConnection_WrongCredentials(t *testing.T) {
	ctx := context.Background()

	startupCtx, startCancel := context.WithTimeout(ctx, containerStartupTimeout)
	defer startCancel()

	container, err := testcontainers.Run(startupCtx, "postgres:16-alpine",
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":       "testuser",
			"POSTGRES_PASSWORD":   "correctpass",
			"POSTGRES_DB":         "testdb",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithCmd("-c", "fsync=off"),
		postgresTmpfsCustomizer{},
	)
	require.NoError(t, err, "failed to start postgres container")

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	hostPort, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	connStr := fmt.Sprintf("postgres://testuser:correctpass@%s:%s/testdb?sslmode=disable", host, hostPort.Port())
	waitErr := waitForPostgresReady(startupCtx, connStr)
	require.NoError(t, waitErr, "postgres did not become ready in time")

	wrongConnStr := fmt.Sprintf("postgres://testuser:wrongpass@%s:%s/testdb?sslmode=disable", host, hostPort.Port())

	client, err := NewConnection(ctx, wrongConnStr)
	assert.Error(t, err, "NewConnection should fail when credentials are rejected by PostgreSQL")
	assert.Nil(t, client, "client should be nil on authentication failure")
}

// --- Client accessor / mutator tests ---

func TestClient_GetConnectionURI(t *testing.T) {
	ctx := context.Background()

	connStr, terminate := newPostgresContainer(ctx, t)
	defer terminate()

	client, err := NewConnection(ctx, connStr)
	require.NoError(t, err)
	defer client.Close()

	assert.NotEmpty(t, client.GetConnectionURI())
}

func TestClient_Close(t *testing.T) {
	ctx := context.Background()

	connStr, terminate := newPostgresContainer(ctx, t)
	defer terminate()

	client, err := NewConnection(ctx, connStr)
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err, "Close should not return an error")
}

func TestClient_SetDB(t *testing.T) {
	ctx := context.Background()

	connStr, terminate := newPostgresContainer(ctx, t)
	defer terminate()

	client, err := NewConnection(ctx, connStr)
	require.NoError(t, err)

	oldDB := client.GetClient()
	assert.NotNil(t, oldDB)

	newDB, err := openSQLxDB(connStr)
	require.NoError(t, err)

	client.SetDB(newDB)
	assert.Same(t, newDB, client.GetClient(), "SetDB should replace the db field")

	err = client.GetClient().PingContext(ctx)
	assert.NoError(t, err)

	_ = newDB.Close()
}

func TestClient_GetClient(t *testing.T) {
	ctx := context.Background()

	connStr, terminate := newPostgresContainer(ctx, t)
	defer terminate()

	client, err := NewConnection(ctx, connStr)
	require.NoError(t, err)

	db := client.GetClient()
	assert.NotNil(t, db)

	err = db.PingContext(ctx)
	assert.NoError(t, err)

	_ = client.Close()
}

func openSQLxDB(connStr string) (*sqlx.DB, error) {
	db, err := sql.Open(string(Postgres), connStr)
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(db, string(Postgres)), nil
}
