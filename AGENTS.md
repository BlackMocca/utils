# AGENTS.md — Project Context for AI Agents

## Repository Overview

- **Name**: `utils`
- **Owner**: BlackMocca
- **URL**: https://github.com/BlackMocca/utils.git
- **Language**: Go (1.26+)
- **Structure**: Multi-module monorepo, 3 independent modules with zero inter-module dependencies
- **Purpose**: Shared utility packages for backend services — struct conversion, typed DB fields, PostgreSQL client

This repository serves as a **knowledge base** for the BlackMocca backend ecosystem's foundational utilities. When making changes, prioritize preserving and enhancing documentation over code-only modifications.

---

## Folder Organization

```
utils/
├── README.md                    ← Entry point: module overview, cross-module relationships
├── AGENTS.md                    ← This file: AI agent context & conventions
├── fn/                          ← Generic struct-to-struct conversion utilities
│   ├── go.mod                   ← Module: github.com/BlackMocca/utils/fn
│   ├── convert.go               ← Core implementation (2 functions, ~30 lines)
│   ├── convert_test.go          ← Tests for both conversion functions
│   └── README.md                ← API reference & usage examples
├── models/                      ← Custom database & JSON serialization types
│   ├── go.mod                   ← Module: github.com/BlackMocca/utils/models
│   ├── constants.go             ← Shared constants (date/time format strings)
│   ├── date.go                  ← Date type (YYYY-MM-DD, no timezone conversion)
│   ├── timestamp.go             ← Timestamp type (full datetime, UTC→Bangkok auto-convert)
│   ├── enum.go                  ← EnumScan generic nullable wrapper
│   ├── json.go                  ← JsonScan generic JSON field wrapper
│   ├── date_test.go             ← Date edge cases: null, empty, invalid formats
│   ├── timestamp_test.go        ← Timestamp format detection & timezone conversion
│   ├── enum_test.go             ← EnumScan nil/zero/JSON behavior
│   └── json_test.go             ← JsonScan maps/arrays/structs + SQL round-trip
└── connectors/
    └── psql/                    ← PostgreSQL connector with pool management
        ├── go.mod               ← Module: github.com/BlackMocca/utils/psql
        ├── client.go            ← Client struct, NewConnection, NewConnectionWithTracing
        ├── client_test.go       ← Integration tests via testcontainers-go
        ├── driver.go            ← Driver type definition (Postgres enum)
        ├── hook.go              ← TracingHook implementation for OpenTracing/sqlhooks
        └── README.md            ← API reference, config, usage examples
```

---

## Module Summaries

### 1. `fn` — Generic Struct & JSON Conversion

**Import path**: `github.com/BlackMocca/utils/fn`  
**Runtime dependencies**: None (standard library only)  
**Test dependencies**: `stretchr/testify`  
**Code size**: ~30 lines (2 functions in one file)

Provides two JSON round-trip conversion functions that map between struct types with different field names — as long as JSON tags align:

- **`ConvertStruct[Src, Dst](src Src) (Dst, error)`** — generic, returns non-pointer destination
- **`CopyJSON(src any, dst any) error`** — pointer-based classic pattern

Both marshal to JSON then unmarshal into the target type. Supports nested structs, slices, maps, nil pointers, and respects `json:"-"`. Not for performance-critical paths or circular references.

**Knowledge context**: This module solves DTO mapping problems in Go where field names differ between layers (API ↔ domain ↔ DB). See [fn README](./fn/README.md) for examples.

---

### 2. `models` — Database & JSON Serialization Types

**Import path**: `github.com/BlackMocca/utils/models`  
**Runtime dependencies**: `4d63.com/tz` (timezone support)  
**Test dependencies**: `stretchr/testify`  
**Files**: 5 source files (`date.go`, `timestamp.go`, `enum.go`, `json.go`, `constants.go`)

Four exported types for database and JSON serialization, all with automatic **Asia/Bangkok** timezone handling:

| Type | File | Purpose | Interfaces Implemented |
|------|------|---------|----------------------|
| `Date` | `date.go` | Date-only (`YYYY-MM-DD`) | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` |
| `Timestamp` | `timestamp.go` | Full datetime with UTC→Bangkok auto-convert | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` |
| `EnumScan[T ~string]` | `enum.go` | Nullable enum wrapper for DB/string types | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` |
| `JsonScan[T any]` | `json.go` | Generic JSON field wrapper (maps, arrays, structs) | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` |

**Critical timezone behavior**:
- `Timestamp.UnmarshalJSON` and `Timestamp.Scan(time.Time)` **automatically convert UTC to Bangkok (UTC+7)**. Local-time values are not converted.
- `Date` does **not** perform timezone conversion — it stores the date as-is.
- This is by design: the team serves a Southeast Asian user base where Bangkok time is the canonical representation.

**Edge case handling**:
- `null`, `""`, and `"nil"` all parse cleanly for `Date` and `Timestamp` without error.
- `EnumScan.MarshalJSON(nil)` returns an empty byte slice; `UnmarshalJSON("")` is a no-op.
- `JsonScan` nil slices marshal to `[]` (empty array), not `null`.

**Knowledge context**: Used everywhere in backend services that interact with PostgreSQL or JSON APIs. See [models README](./models/README.md) for full API reference.

---

### 3. `psql` — PostgreSQL Connector

**Import path**: `github.com/BlackMocca/utils/psql`  
**Dependencies**: pgx/v5, sqlx, lib/pq (fallback), opentracing-go, sqlhooks/v2, spf13/cast  
**Test dependencies**: testcontainers-go (Docker integration tests)

Single `Client` struct wrapping a pooled `*sqlx.DB`:

- **`NewConnection(ctx, dsn)`** — default client
- **`NewConnectionWithTracing(ctx, dsn, driver, tracer)`** — with OpenTracing hooks on SELECT/INSERT/UPDATE/DELETE

Pool settings configurable via environment variables (`POSTGRES_*`). Connection pooling uses `pgxpool`.

**Knowledge context**: The database access layer used across all backend services. Integrates with `models` types for typed columns. See [psql README](./connectors/psql/README.md) for full reference.

---

## Running Tests

```bash
# fn module (pure unit tests, no Docker)
cd fn && go test ./... -v

# models module (pure unit tests, no Docker)
cd models && go test ./... -v

# psql module (integration tests require Docker via testcontainers-go)
cd connectors/psql && go test -v -count=1 ./...
```

The `psql` tests use `testcontainers-go` to spin up a temporary PostgreSQL 16 container with tmpfs-mounted data directory (avoids WSL2 overlayfs issues).

---

## Documentation Conventions

### What constitutes source of truth

1. **Source code** — always the primary reference for API signatures, behavior, and edge cases
2. **Module README.md** — detailed usage examples, API reference, gotchas
3. **This AGENTS.md** — cross-module context, knowledge graph awareness, agent workflow guidance
4. **README.md (root)** — entry point for humans and agents navigating the repo

### File naming and structure rules

- Module-level documentation lives in `module/README.md` (not per-file docs)
- This file is at module root as `AGENTS.md`
- No ADR directory yet — document design decisions inline in relevant module READMEs or AGENTS.md sections
- Do **not** create duplicate documentation. If information exists in one place, reference it rather than repeating it.

### Cross-referencing policy

When writing new documentation or modifying existing code:
- **Always link to related concepts** in other modules using relative paths (e.g., `[→ models documentation](../models/README.md)`)
- When a function in `fn` is commonly used with types from `models`, mention that pairing
- When a type in `models` is tested against SQL operations, note its compatibility with `psql.Client`
- Preserve rationale — if you add a new type, document *why* it exists alongside *what* it does

---

## Knowledge Graph Considerations

This repository is designed to be indexable by knowledge graph tools (graphify, vector databases, RAG systems). When working here:

### Entity extraction priorities
- **High priority**: README files, AGENTS.md, source code (`.go`), module structure (`go.mod`)
- **Low priority / ignore**: test files, coverage reports, build artifacts, `.git` internals

### Relationship mapping guidance
When analyzing this repository, identify and preserve these relationship types:
- **Structural**: `module A → import path of module B` (none currently — modules are independent)
- **Semantic**: `fn.ConvertStruct` ↔ `models.Timestamp` (commonly paired for DTO mapping)
- **Usage**: `psql.Client` + `models.Date/Timestamp` (typed columns in queries)
- **Conceptual**: `Date` vs `Timestamp` (same concept, different granularity — both implement same interfaces)

### Node naming conventions (for graphify / knowledge indexing)
- Use lowercase with underscores: `date_type`, `timestamp_serialization`, `psql_pool_config`
- One level of parent only: `{filename}_{entity}` — e.g., `timestamp_timestamp` for the Timestamp type in timestamp.go
- Do not append chunk or sequence suffixes

---

## Agent Workflow Guidelines

### When making changes

1. **Read before writing** — check if documentation already covers the topic you're adding
2. **Prefer updating over creating** — extend existing README sections rather than creating new files, unless there's clear separation of concern
3. **Add cross-references** — when modifying `models/README.md` to add a type, note how it integrates with `psql` in both docs
4. **Preserve rationale** — if you change behavior (e.g., timezone handling), update the "Why" section alongside the "How"

### When adding new modules

1. Create a top-level entry in root `README.md` (Modules table)
2. Add an `AGENTS.md` section with the same structure as existing modules
3. Create `module/README.md` following the pattern: Install → Overview → Types/Functions → Examples → Testing
4. Ensure the new module's `go.mod` is properly structured

### When fixing documentation bugs

- Update the source of truth (usually README.md), not just this AGENTS.md summary
- If the fix changes behavior, note it in an "Edge Cases" or "Gotchas" section
- Keep the AGENTS.md summary aligned with actual code behavior — it's your agent context contract

---

## Architecture Decisions & Rationale

### Why multi-module monorepo?
Each module is independently deployable as a Go package. Independent `go.mod` files allow different version numbers, dependencies, and release cycles without affecting other modules.

### Why JSON round-trip for struct conversion (fn)?
Simplicity over performance — the functions are 30 lines total vs hundreds of reflection-based alternatives. Acceptable for DTO mapping where correctness matters more than micro-optimization. See `fn/README.md` "When to Use" section for limitations.

### Why Asia/Bangkok timezone in models?
The primary user base operates in Indochina Time (ICT, UTC+7). Storing and displaying in Bangkok time avoids timezone confusion at the application layer. Timestamp auto-conversion from UTC on unmarshal/scan ensures external data is normalized to the canonical timezone.

### Why no inter-module dependencies?
Each module can be imported independently by services that only need one capability. Adding a dependency on `models` would pull in `4d63.com/tz` for all consumers; adding `fn` would add zero overhead. This independence has proven valuable — some services only use `fn`, others only use `psql`.

### Why psql wraps sqlx instead of using pgx directly?
`sqlx` provides convenient struct scanning, named parameters, and a familiar API. The `pgx` driver is used underneath (`sqlx.Open("pgx", dsn)`), giving us both performance and ergonomics.

---

## License

Internal utility packages — see repository root for license details.
