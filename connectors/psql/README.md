# psql — PostgreSQL Connector for Go

A lightweight, production-ready PostgreSQL connection wrapper built on top of [`pgx/v5`](https://github.com/jackc/pgx) and [`sqlx`](https://github.com/jmoiron/sqlx). Provides a unified `Client` struct with pool management, connection lifecycle handling, and optional OpenTracing support.

> **Part of the [utils](../../README.md) multi-module monorepo.** See [AGENTS.md](../../AGENTS.md) for cross-module relationship context.

---

## Installation

```bash
go get github.com/BlackMocca/utils/psql@latest
```

### Dependencies (automatically fetched)

| Package | Purpose |
|---------|---------|
| `github.com/jackc/pgx/v5` | PostgreSQL driver & connection pool (`pgxpool`) |
| `github.com/jmoiron/sqlx` | Extended SQL operations |
| `github.com/lib/pq` | Fallback driver registration |
| `github.com/opentracing/opentracing-go` | Distributed tracing hooks |
| `github.com/qustavo/sqlhooks/v2` | Query lifecycle hooks (`Before`, `After`, `OnError`) |
| `github.com/spf13/cast` | Argument serialization for spans |

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    psql "github.com/BlackMocca/utils/psql"
)

func main() {
    ctx := context.Background()

    // 1. Create a connection using a PostgreSQL DSN
    connStr := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
    
    client, err := psql.NewConnection(ctx, connStr)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // 2. Execute queries using the embedded sqlx.DB
    var count int
    err = client.GetClient().GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Total users: %d\n", count)
}
```

---

## Usage

### 1. Basic Connection (No Tracing)

```go
import psql "github.com/BlackMocca/utils/psql"

// Create client with default settings
client, err := psql.NewConnection(ctx, dsn)
if err != nil {
    // Handle connection error
}
defer client.Close()

// Use sqlx methods directly
err = client.GetClient().PingContext(ctx)
```

### 2. Connection with OpenTracing

When you need distributed tracing (e.g., Jaeger, Zipkin):

```go
import (
    "github.com/opentracing/opentracing-go"
    psql "github.com/BlackMocca/utils/psql"
)

// Initialize your tracer (example: Jaeger)
tracer := // ... initialize opentracing.Tracer instance

client, err := psql.NewConnectionWithTracing(
    ctx, 
    dsn,
    psql.Postgres,
    tracer,
)
if err != nil {
    panic(err)
}
defer client.Close()
```

### 3. Replace the Underlying DB Instance

Useful for migrations, test fixtures, or connection pooling changes:

```go
newDB, _ := sqlx.Open("pgx", newDSN)
client.SetDB(newDB) // Replaces internal *sqlx.DB
```

---

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `NewConnection(ctx, dsn)` | Creates a PostgreSQL client. Pool settings (`MaxConns`, `MaxConnIdleTime`, `MaxConnLifetime`) are read from `POSTGRES_*` env vars or fall back to defaults (4, 30m, 1h). Returns `*Client` or error. |
| `NewConnectionWithTracing(ctx, dsn, driver, tracer)` | Same as above but registers a tracing hook that emits spans for SELECT/INSERT/UPDATE/DELETE queries. |

### Methods on `Client`

| Method | Description |
|--------|-------------|
| `GetClient() *sqlx.DB` | Returns the underlying `*sqlx.DB` instance (pool + connection). |
| `GetConnectionURI() string` | Returns the DSN used to create the client. |
| `SetDB(db *sqlx.DB)` | Replaces the current database connection with a new one. |
| `Close() error` | Closes all pooled connections and releases resources. |

### Types

#### `Driver` (String Enum)

```go
type Driver string

const Postgres Driver = "pgx"  // Currently supported driver
```

Use `psql.Postgres` when calling `NewConnectionWithTracing`.

---

## Configuration Defaults

| Setting | Value |
|---------|-------|
| Max Connections | 4 |
| Max Idle Time | 30 minutes |
| Max Lifetime | 1 hour |
| Pool Mode | `pgxpool` (managed by jackc/pgx) |
| SSL | Negotiated via DSN parameters (`sslmode=disable` recommended for local/dev) |

### Environment Variables

All pool settings can be overridden via environment variables with the `POSTGRES_` prefix. If a variable is missing or contains an invalid value, the default is used.

| Variable | Default | Format / Example |
|----------|---------|------------------|
| `POSTGRES_MAX_CONNS` | `4` | Integer > 0 (e.g., `16`) |
| `POSTGRES_MAX_CONN_IDLE_TIME` | `30m` | Go duration string (e.g., `5m`, `1h`, `3600s`) |
| `POSTGRES_MAX_CONN_LIFETIME` | `1h` | Go duration string (e.g., `30m`, `2h`, `7200s`) |

**Examples:**

```bash
# Override all pool settings for a production service
export POSTGRES_MAX_CONNS=20
export POSTGRES_MAX_CONN_IDLE_TIME=5m
export POSTGRES_MAX_CONN_LIFETIME=45m
```

Values are read at connection time — changing env vars after `NewConnection` has been called has no effect (you must create a new client).

---

## Connection String Examples

### Local Development

```
postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Production (TLS)

```
postgres://user:password@prod-host:5432/dbname?sslmode=require&sslrootcert=/path/to/ca.pem
```

### Docker / Internal Network

```
postgres://app_user:secret@psql-container:5432/app_db?sslmode=disable
```

---

## Testing

Unit tests use `testcontainers-go` to spin up a temporary PostgreSQL 16 container with tmpfs-mounted data directory (avoids overlayfs permission issues in WSL2/CI environments).

```bash
cd /path/to/module

# Run all tests
go test -v -count=1 ./...

# Run single test
go test -v -run TestNewConnection_Success ./...

# Verify code quality
go vet ./...
```

### Requirements for Tests

- **Docker** running locally (for integration tests)
- Go 1.26+ module mode enabled
- Dependencies in `go.mod`:
  ```
  testcontainers-go
  github.com/jackc/pgx/v5/stdlib
  github.com/stretchr/testify
  ```

---

## Error Handling

All public functions return `(result, error)` pairs. Common errors include:

| Scenario | Error Type |
|----------|-----------|
| Invalid DSN / malformed URL | `error` from pgx pool creation |
| Network unreachable / port closed | Connection refused / timeout |
| Authentication failed (wrong credentials) | PostgreSQL auth error (`FATAL: password authentication failed`) |
| SSL negotiation failure | TLS handshake error |

Always check the returned `err` before using the client.

---

## Related Modules

### With [models](../../models/README.md)
The psql client works seamlessly with `models` types for typed database columns:

```go
import "github.com/BlackMocca/utils/models"

type User struct {
    ID        uint               
    CreatedAt models.Timestamp   // Auto-converts UTC → Bangkok on scan
    Role      models.EnumScan[RoleType]
}

var user User
client.GetClient().GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", 1)
// user.CreatedAt is already in Asia/Bangkok timezone
```

### With [fn](../../fn/README.md)
Use `ConvertStruct` to map between database models (containing `models.*` types) and API DTOs:

```go
import "github.com/BlackMocca/utils/fn"

type UserAPI struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}

dst, err := fn.ConvertStruct[User, UserAPI](user)
```

---

## License

Internal utility package — see repository root for license details.
