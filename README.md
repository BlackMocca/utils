# utils — Shared Go Utility Packages

[![Go Reference](https://pkg.go.dev/badge/github.com/BlackMocca/utils.svg)](https://pkg.go.dev/github.com/BlackMocca/utils)

Lightweight, production-ready Go utility packages used across the BlackMocca backend ecosystem. This is a **multi-module monorepo** with three independent modules that share no internal dependencies — each can be imported standalone or together as needed.

```
utils/
├── README.md                 ← You are here: overview & navigation
├── AGENTS.md                 ← AI agent context & conventions
├── fn/                       ← Generic struct-to-struct conversion
├── models/                   ← Custom DB & JSON serialization types
└── connectors/psql/          ← PostgreSQL client with pool management
```

---

## Modules

| Module | Import Path | Purpose | Key Types / Functions |
|--------|-------------|---------|----------------------|
| **[fn](./fn)** | `github.com/BlackMocca/utils/fn` | Struct-to-struct conversion via JSON round-trip | `ConvertStruct[T, U]`, `CopyJSON` |
| **[models](./models)** | `github.com/BlackMocca/utils/models` | Database & JSON types: Date, Timestamp, EnumScan, JsonScan | `Date`, `Timestamp`, `EnumScan[T]`, `JsonScan[T]` |
| **[psql](./connectors/psql)** | `github.com/BlackMocca/utils/connectors/psql` | PostgreSQL client wrapper over pgx/v5 + sqlx with pool management & optional tracing | `Client`, `NewConnection`, `NewConnectionWithTracing` |

All modules require **Go 1.26+**. They share no dependencies between each other — `fn` and `models` have zero runtime deps beyond the standard library; `psql` adds database drivers, sqlx, tracing, and test tooling.

---

## Quick Install

```bash
# Generic conversion utilities (no runtime deps)
go get github.com/BlackMocca/utils/fn@latest

# Custom date/time, enum, and JSON types for DB/JSON serialization
go get github.com/BlackMocca/utils/models@latest

# PostgreSQL connector wrapper
go get github.com/BlackMocca/utils/connectors/psql@latest
```

---

## Cross-Module Relationships

While each module is independently importable, they are commonly used together in backend services:

### models + psql (most common pairing)

`models` types implement `database/sql.Scanner` and `driver.Valuer`, making them directly usable as struct fields with the `psql` client:

```go
import "github.com/BlackMocca/utils/models"

type User struct {
    ID        uint               `gorm:"primary_key"`
    CreatedAt models.Timestamp   `gorm:"column:created_at;type:timestamp"`
    Status    models.EnumScan[string]
}
// Works with psql.Client.GetClient().GetContext(...) or GORM
```

[→ models documentation](./models/README.md) · [→ psql documentation](./connectors/psql/README.md)

### fn + models (DTO transformation)

`fn.ConvertStruct` is useful for mapping between `models` types and API DTOs when field names differ:

```go
import "github.com/BlackMocca/utils/fn"
import "github.com/BlackMocca/utils/models"

// Map internal timestamp to API-friendly string representation
apiResponse, err := fn.ConvertStruct[models.Timestamp, map[string]any](ts)
```

[→ fn documentation](./fn/README.md)

### All three (full service stack)

A typical BlackMocca backend service imports all three: `fn` for DTOs, `models` for typed database fields, and `psql` for the connection layer.

---

## Development Conventions

- Each module is an **independent** Go module with its own `go.mod` / `go.sum`.
- All modules target **Go 1.26+**.
- Tests use `stretchr/testify` for assertions.
- No global mutable state in `fn`; all functions are pure.
- `models` types implement standard interfaces (`Scanner`, `Valuer`, `Marshaler`, `Unmarshaler`).
- The `psql` module uses `sqlhooks/v2` for query lifecycle hooks and `opentracing-go` for distributed tracing.

---

## Running Tests

```bash
# fn module (no Docker needed)
cd fn && go test ./... -v

# models module (no Docker needed)
cd models && go test ./... -v

# psql module (requires Docker for integration tests via testcontainers-go)
cd connectors/psql && go test -v -count=1 ./...
```

---

## License

Internal utility packages — see repository root for license details.
