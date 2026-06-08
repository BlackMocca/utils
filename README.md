# utils — Shared Go Utility Packages

[![Go Reference](https://pkg.go.dev/badge/github.com/BlackMocca/utils.svg)](https://pkg.go.dev/github.com/BlackMocca/utils)

A collection of lightweight, production-ready Go utility packages used across the BlackMocca backend ecosystem. This is a multi-module monorepo with three independent modules:

| Module | Import Path | Description |
|--------|-------------|-------------|
| **[fn](./fn)** | `github.com/BlackMocca/utils/fn` | Generic struct-to-struct conversion via JSON round-trip |
| **[models](./models)** | `github.com/BlackMocca/utils/models` | Custom DB & JSON types: Date, Timestamp, EnumScan, JsonScan |
| **[psql](./connectors/psql)** | `github.com/BlackMocca/utils/psql` | PostgreSQL client with pool management & optional tracing |

All modules require **Go 1.26+**.

---

## Quick Install

```bash
# Generic conversion utilities
go get github.com/BlackMocca/utils/fn@latest

# Custom date/time, enum, and JSON types for DB/JSON serialization
go get github.com/BlackMocca/utils/models@latest

# PostgreSQL connector wrapper
go get github.com/BlackMocca/utils/psql@latest
```

---

## 1. fn — Struct Conversion Utilities

Lightweight functions to copy data between structs via `encoding/json` round-trip. Field names don't need to match — only JSON tags matter. Supports nested structs, slices, maps, pointers (including nil), and respects `json:"-"`.

```go
import "github.com/BlackMocca/utils/fn"

type Source struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
}

type Dest struct {
    Identifier int    `json:"id"`
    FullName   string `json:"name"`
}

// Generic: caller specifies types at compile time
dst, err := fn.ConvertStruct[Source, Dest](src)
if err != nil { /* handle */ }

// Pointer-based (classic pattern)
var dst Dest
err = fn.CopyJSON(src, &dst)
```

[→ Full documentation](./fn/README.md)

---

## 2. models — Date, Timestamp, Enum & JSON Types

Custom types for database and JSON serialization with automatic **Asia/Bangkok** (UTC+7) timezone handling. Designed for use with GORM / raw `database/sql`.

```go
import "github.com/BlackMocca/utils/models"

// Date (YYYY-MM-DD only)
d := models.NewDateFromString("2024-06-15")

// Timestamp — auto-converts UTC → Bangkok on unmarshal/scan
ts := models.NewTimestampFromNow()

// Nullable enum for DB columns
type Status string // "active" | "inactive"
e := models.NewEnumScan[Status]("active")

// Generic JSON field (maps, arrays, structs)
js := models.NewJsonScan(map[string]any{"color": "red"})
```

All types implement `Scanner`, `Valuer`, `Marshaler`, and `Unmarshaler` interfaces.

[→ Full documentation](./models/README.md)

---

## 3. psql — PostgreSQL Connector

A production-ready PostgreSQL client built on `pgx/v5` + `sqlx`. Provides pool management, configurable via environment variables, with optional OpenTracing support.

```go
import "github.com/BlackMocca/utils/psql"

ctx := context.Background()
client, err := psql.NewConnection(ctx, dsn)
if err != nil { panic(err) }
defer client.Close()

// Use the underlying *sqlx.DB directly
var count int
err = client.GetClient().GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
```

**Pool settings** are configured via `POSTGRES_*` env vars:

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_MAX_CONNS` | `4` | Max open connections |
| `POSTGRES_MAX_CONN_IDLE_TIME` | `30m` | Idle timeout |
| `POSTGRES_MAX_CONN_LIFETIME` | `1h` | Connection max age |

For distributed tracing, use `NewConnectionWithTracing(ctx, dsn, psql.Postgres, tracer)`.

[→ Full documentation](./connectors/psql/README.md)

---

## Running Tests

Each module has its own test suite. Run from the module directory:

```bash
# fn module
cd fn && go test ./... -v

# models module
cd models && go test ./... -v

# psql module (requires Docker for integration tests)
cd connectors/psql && go test -v -count=1 ./...
```

---

## License

Internal utility packages — see repository root for license details.
