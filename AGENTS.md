# AGENTS.md — Project Context for AI Agents

## Repository

- **Name**: `utils`
- **Owner**: BlackMocca
- **URL**: https://github.com/BlackMocca/utils.git
- **Language**: Go (1.26+)
- **Structure**: Multi-module monorepo (3 independent modules)

---

## Module Overview

| Module | Import Path | Description |
|--------|-------------|-------------|
| `fn` | `github.com/BlackMocca/utils/fn` | Generic struct-to-struct conversion via JSON round-trip |
| `models` | `github.com/BlackMocca/utils/models` | Custom DB/JSON types: Date, Timestamp, EnumScan, JsonScan (Asia/Bangkok timezone) |
| `psql` | `github.com/BlackMocca/utils/psql` | PostgreSQL client wrapper over pgx/v5 + sqlx with pool management & optional tracing |

---

## Module Details

### 1. `fn` — Generic Struct & JSON Conversion

**Package**: `github.com/BlackMocca/utils/fn`  
**Go Version**: 1.26+  
**Dependencies**: Only `testify` (tests)

Exports two functions that marshal to JSON and unmarshal into another type:

- **`ConvertStruct[Src, Dst](src Src) (Dst, error)`** — generic function, caller specifies both types at compile time, returns a non-pointer value.
- **`CopyJSON(src any, dst any) error`** — classic pointer-based pattern.

Both use `encoding/json` under the hood. Field names don't need to match; only JSON tags must align. Supports nested structs, slices, maps, pointers (including nil), and respects `json:"-"` exclusion.

Not for performance-critical paths or circular references.

---

### 2. `models` — Database & Serialization Types

**Package**: `github.com/BlackMocca/utils/models`  
**Go Version**: 1.26+  
**Dependencies**: `4d63.com/tz` (timezone support)

All types store data in **Asia/Bangkok** timezone (UTC+7). Four exported types:

| Type | Purpose | Key Methods / Constructors |
|------|---------|----------------------------|
| `Date` | Date-only (`YYYY-MM-DD`) wrapped from `time.Time` | `NewDateFromString`, `ParseDateFromString`, `Format`, `Weekday`, `ToTime`, `ToTimestamp`, `Scan`, `Value` |
| `Timestamp` | Full timestamp with **UTC → Bangkok auto-conversion** | `NewTimestampFromNow`, `NewTimestampFromString`, `ParseTimestampFromString`, `ToUnix`, `YearDay`, `ValueOrZero`, `Scan`, `Value` |
| `EnumScan[T ~string]` | Nullable enum wrapper for DB/string types | `NewEnumScan`, `Data`, `Set`, `String`, `IsZero`, `Scan`, `Value` |
| `JsonScan[T any]` | Generic JSON field wrapper (maps, arrays, structs) | `NewJsonScan`, `Data`, `Set` (chainable), `Scan`, `Value` |

**Important behavior notes**:
- `Timestamp.UnmarshalJSON` auto-detects 14+ formats including RFC3339, microseconds/milliseconds variants. UTC timestamps are converted to Bangkok (UTC+7).
- `Date` and `Timestamp` unmarshal handle `null`, empty string (`""`), and `"nil"` without error.
- `EnumScan.MarshalJSON(nil)` returns an empty byte slice; `UnmarshalJSON("")` is a no-op.
- `JsonScan` nil slices marshal to `[]` (empty array), not `null`.
- `Date.Value()` and `Timestamp.Value()` return `nil` for zero values — useful for nullable DB columns.

---

### 3. `psql` — PostgreSQL Connector

**Package**: `github.com/BlackMocca/utils/psql`  
**Go Version**: 1.26+  
**Dependencies**: pgx/v5, sqlx, lib/pq, opentracing-go, sqlhooks/v2, spf13/cast, testcontainers (tests)

Provides a single `Client` struct wrapping a pooled `*sqlx.DB`:

| Function | Description |
|----------|-------------|
| `NewConnection(ctx, dsn) (*Client, error)` | Default client (no tracing) |
| `NewConnectionWithTracing(ctx, dsn, driver, tracer)` | Client with OpenTracing hooks on SELECT/INSERT/UPDATE/DELETE |

**Methods on `Client`**:
- `GetClient() *sqlx.DB` — access the underlying db
- `GetConnectionURI() string` — returns the DSN
- `SetDB(db *sqlx.DB)` — swap the internal db instance
- `Close() error` — close pool

**Pool defaults** (overridable via env vars):

| Setting | Env Variable | Default |
|---------|-------------|---------|
| Max connections | `POSTGRES_MAX_CONNS` | 4 |
| Max idle time | `POSTGRES_MAX_CONN_IDLE_TIME` | 30m |
| Max lifetime | `POSTGRES_MAX_CONN_LIFETIME` | 1h |

---

## Running Tests

All modules live under the same repository root. Run tests per module:

```bash
# fn module
cd fn && go test ./... -v

# models module
cd models && go test ./... -v

# psql module (requires Docker for integration tests)
cd connectors/psql && go test -v -count=1 ./...
```

The `psql` tests use `testcontainers-go` to spin up a temporary PostgreSQL 16 container with tmpfs-mounted data directory (avoids WSL2 overlayfs issues).

---

## Development Conventions

- Each module is an **independent** Go module with its own `go.mod` / `go.sum`.
- All modules target **Go 1.26+**.
- Tests use `stretchr/testify` for assertions.
- The `psql` module uses `sqlhooks/v2` for query lifecycle hooks and `opentracing-go` for distributed tracing integration.
- No global mutable state in `fn`; pure functions only.
- `models` types implement `database/sql.Scanner`, `driver.Valuer`, `json.Marshaler`, and `json.Unmarshaler` interfaces where applicable.
