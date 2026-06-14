# AGENTS.md — models Package Knowledge Base for AI Agents

## Repository Purpose

This knowledge base documents the **BlackMocca `utils/models`** Go package — a collection of custom types for database and JSON serialization with automatic **Asia/Bangkok timezone handling**. It is designed to be read by AI agents before they modify, extend, or integrate this package.

### When to Read This File

| Situation | Why Read AGENTS.md First |
|-----------|-------------------------|
| Modifying a type (`Date`, `Timestamp`, `EnumScan`, `JsonScan`) | Understand timezone rules, JSON/SQL contracts, edge cases before editing source |
| Adding a new type or method | Follow naming conventions, interface implementations (Scanner/Valuer/Marshaler) |
| Writing tests | Know which edge cases are already covered in `_test.go` files and what's missing |

### Document Map

| Document | Purpose | When to Read |
|----------|---------|--------------|
| **This file** (AGENTS.md) | Agent context, timezone rules, type contracts, integration patterns | Starting any task involving the `models` package |
| [README.md](./README.md) | API reference: constructors, methods, JSON/SQL examples for each type | Looking up a specific signature or usage pattern |
| Source files (`date.go`, `timestamp.go`, `enum.go`, `json.go`) | Implementation details — exact parsing logic, timezone conversion code | Debugging an edge case or verifying behavior |

---

## Package Overview

**Package**: `github.com/BlackMocca/utils/models`  
**Go Version**: 1.26+ (uses generics for `EnumScan[T]`, `JsonScan[T]`)  
**Dependencies**: `4d63.com/tz` (timezone support — Asia/Bangkok)

### The Four Types

| Type | Wraps | Purpose | Implements | Source File |
|------|-------|---------|------------|-------------|
| **Date** | `time.Time` | Date-only operations (`YYYY-MM-DD`) | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` | `date.go` |
| **Timestamp** | `time.Time` | Full timestamp with UTC→Bangkok auto-conversion | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` | `timestamp.go` |
| **EnumScan[T]** | `*T` (pointer to string type) | Nullable enum/string with DB & JSON support | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` | `enum.go` |
| **JsonScan[T]** | `*T` (pointer, any type) | Generic JSON column — maps, arrays, or typed structs | `Scanner`, `Valuer`, `Marshaler`, `Unmarshaler` | `json.go` |

---

## Core Concepts & Timezone Rules

### The One Rule Everyone Must Know: UTC→Bangkok Conversion on Timestamp

`Timestamp` automatically converts **UTC → Asia/Bangkok (ICT, UTC+7)** during both JSON unmarshaling and SQL scanning. Local-timezone values are NOT converted.

```
Incoming value with location == time.UTC
  → subtract 7 hours (t.Add(-7 * time.Hour))
  → store as Timestamp in Location (Asia/Bangkok)

Incoming value with any other location
  → store as-is, no conversion

Marshal/Value output
  → always uses "2006-01-02 15:04:05" format — displayed in Bangkok time
```

**Why this matters for agents**: When reading a timestamp from an external API or database, the value you see is **already converted to UTC+7**. If the input was UTC "2025-04-03 12:05:12Z", the stored and displayed value will be "2025-04-03 19:05:12".

### Date Type — No Timezone Conversion

`Date` does **not** perform timezone conversion. It stores whatever is passed, parsing only the date portion (`YYYY-MM-DD`). Use `Date` when timezone granularity doesn't matter (birth dates, scheduled dates).

---

## Type Contracts Reference

### Date Contract

| Operation | Input → Output | Notes |
|-----------|---------------|-------|
| `NewDateFromString("2024-06-15")` | `Date` | Panics on invalid input — use `ParseDateFromString` for error handling |
| `ParseDateFromString("2024-06-15")` | `(Date, error)` | Safe alternative |
| `NewDateFromTime(t time.Time)` | `Date` | Extracts date portion only |
| `d.Format(layout)` | `string` | Any Go time layout — e.g., `"Mon Jan 2 2006"` |
| `d.String()` | `string` | Always `"YYYY-MM-DD"` |
| `d.ToTime()` | `time.Time` | Convert back to standard `time.Time` |
| `d.ToTimestamp()` | `Timestamp` | Promote date → full timestamp (time = 00:00:00 Bangkok) |
| `d.Value()` | `(driver.Value, error)` | Returns `"YYYY-MM-DD"` or `nil` if zero date |
| `json.Marshal(d)` | `[]byte` | Always outputs `"2024-06-15"` format |

**Null/empty handling**: JSON unmarshal ignores `null`, `""`, `"nil"` — no change to the field. Invalid strings return error.

### Timestamp Contract

| Operation | Input → Output | Notes |
|-----------|---------------|-------|
| `NewTimestampFromNow()` | `Timestamp` | Current time in Bangkok |
| `NewTimestampFromString("2024-06-15 14:30:45")` | `Timestamp` | Panics on invalid input — use `ParseTimestampFromString` for error handling |
| `ParseTimestampFromString(s)` | `(Timestamp, error)` | Safe alternative; empty string → zero timestamp without error |
| `ts.String()` | `string` | Always `"YYYY-MM-DD HH:MM:SS"` (Bangkok time) |
| `ts.ToTime()` | `time.Time` | Convert back to standard `time.Time` |
| `ts.ToUnix()` | `int64` | Unix seconds — uses internal `time.Time.Unix()` which is always UTC-based |
| `ts.ValueOrZero()` | `string` | Returns formatted string or `""` if zero (safe for nullable fields) |
| `ts.Value()` | `(driver.Value, error)` | Returns `"YYYY-MM-DD HH:MM:SS"` or `nil` if zero |
| `json.Marshal(ts)` | `[]byte` | Always outputs `"2024-06-15 14:30:45"` format (Bangkok time) |

**Supported JSON input formats**: Standard, microseconds (`999999`), milliseconds (`999`), ISO8601 (`T` separator), RFC3339, RFC3339Nano, RFC1123, UnixDate, ANSIC, and date-only.

### EnumScan[T] Contract

| Operation | Input → Output | Notes |
|-----------|---------------|-------|
| `NewEnumScan[string]("active")` | `EnumScan[string]` | Creates non-nil enum scan from string |
| `e.Data()` | `T` | Returns zero value if nil; never returns nil pointer |
| `e.Set(v T)` | `void` | Sets new value |
| `e.String()` | `string` | `""` if nil, else the value |
| `e.IsZero()` | `bool` | `true` if v is nil or empty string |
| `e.Value()` | `(driver.Value, error)` | Returns string or `nil` if zero |

**Null handling**: SQL `Scan(nil)` → sets to zero (no error). JSON `Marshal(nil)` → returns `[]byte("")` (empty slice). JSON unmarshal ignores empty input.

### JsonScan[T] Contract

| Operation | Input → Output | Notes |
|-----------|---------------|-------|
| `NewJsonScan(data)` | `JsonScan[T]` | Creates from **any** value — maps, arrays, typed structs, pointer slices |
| `js.Data()` | `T` | Allocates zero if nil — never returns nil pointer |
| `js.Set(v T)` | `*JsonScan[T]` | Fluent/chaining — returns pointer for builder pattern |
| `js.Value()` | `(driver.Value, error)` | Returns JSON-encoded `[]byte` or `nil` if zero |

**Supported `T` types**: `map[string]interface{}`, `[]interface{}`, typed structs (`Profile`, `Config`), slice of pointers (`[]*Permission`, `[]*Role`), nested structs.

**Null/empty handling**: SQL `Scan(nil)` → no change. JSON unmarshal ignores empty input. **Special rule**: nil slices marshal to `[]` (empty array), not `null`.

---

## Common Usage Patterns

### Pattern 1: Domain Model with All Four Types

```go
type User struct {
    ID          uint               
    CreatedAt   models.Timestamp              // Auto-UTC→Bangkok conversion on scan
    BirthDate   models.Date                   // Date only, no timezone issues  
    Role        models.EnumScan[constants.Role]  // Nullable enum from DB
    Profile     models.JsonScan[map[string]interface{}]  // Flexible JSON column
}
```

### Pattern 2: Nullable Fields with Pointers

```go
type Article struct {
    PublishedAt *models.Timestamp  // nil = not yet published
    EffectiveDate *models.Date     // nil = no date restriction
}
```

### Pattern 3: Fluent JSON Field Updates (maps/arrays)

```go
config := NewJsonScan(map[string]interface{}{
    "color": "red",
    "size":  "L",
})

// Update via fluent Set (returns *JsonScan for chaining)
config.Set(map[string]interface{}{"color": "blue"})

// Retrieve — Data() allocates zero if nil, never panics
data := config.Data()
```

### Pattern 4: JsonScan with Struct Pointers (`[]*Permission`)

Use `JsonScan[T]` to store **structured data** (not just raw maps) in a single JSON column:

```go
// Domain struct definition
type Permission struct {
    ID        uuid.UUID          `json:"id" db:"id"`
    Name      string             `json:"name" db:"name"`
    CreatedAt *models.Timestamp  `json:"created_at" db:"created_at"`
}

// User model with JsonScan holding a slice of *Permission
type User struct {
    ID          uint                        `json:"id" db:"id"`
    Name        string                      `json:"name" db:"name"`
    Role        models.EnumScan[constants.Role] `json:"role" db:"role"`
    Permissions models.JsonScan[[]*Permission]  `json:"permissions" db:"permissions"`
}
```

#### CRUD Example — Permissions

```go
// Create user with permissions
dbPermissions := []*Permission{
    {ID: uuid.New(), Name: "read", CreatedAt: models.NewTimestampFromNow().ToPointer()},
    {ID: uuid.New(), Name: "write", CreatedAt: models.NewTimestampFromNow().ToPointer()},
}

user := User{
    ID:          1,
    Name:        "Alice",
    Role:        NewEnumScan[constants.Role](string(constants.RoleAdmin)),
    Permissions: *NewJsonScan(dbPermissions),  // store as JsonScan
}

// Save to DB (permissions → JSON column)
db.Create(&user)

// Read from DB — permissions auto-unmarshaled back to []*Permission
db.First(&user, 1)
perms := user.Permissions.Data()  // returns []*Permission
for _, perm := range perms {
    fmt.Printf("%s: %s\n", perm.ID, perm.Name)  // "uuid-...: read"
}

// Update permissions
newPerms := []*Permission{
    {ID: uuid.New(), Name: "delete"},
}
user.Permissions.Set(newPerms)
db.Save(&user)
```

### Pattern 5: JsonScan with Single Struct (`Profile`)

Use `JsonScan[T]` for a **single nested struct** instead of multiple JSON columns:

```go
type Profile struct {
    Bio      string         `json:"bio" db:"bio"`
    Avatar   string         `json:"avatar" db:"avatar"`
    Settings map[string]interface{} `json:"settings" db:"settings"`
}

type User struct {
    ID       uint                  `json:"id" db:"id"`
    Name     string                `json:"name" db:"name"`
    Profile  models.JsonScan[Profile]  `json:"profile" db:"profile"`
}
```

#### CRUD Example — Single Struct

```go
// Create user with profile
user := User{
    ID: 1,
    Name: "Bob",
    Profile: *NewJsonScan(Profile{
        Bio:      "Software engineer",
        Avatar:   "https://example.com/bob.png",
        Settings: map[string]interface{}{"theme": "dark", "lang": "th"},
    }),
}
db.Create(&user)

// Read and modify profile
var fetchedUser User
db.First(&fetchedUser, 1)
profile := fetchedUser.Profile.Data()  // returns Profile (never nil)
profile.Bio = "Senior engineer"
fetchedUser.Profile.Set(profile)      // fluent Set for chaining
db.Save(&fetchedUser)
```

### Pattern 6: Safe Zero Value Access

```go
ts := models.Timestamp{}  // zero timestamp
fmt.Println(ts.ValueOrZero())  // "" (no panic, no nil pointer)

e := EnumScan[string]{}
fmt.Println(e.IsZero())    // true
fmt.Println(e.Data())      // "" (zero value, never panics)
```

---

## Integration with psql Client

All four types implement `database/sql.Scanner` and `driver.Valuer`, enabling direct use as struct fields without manual conversion:

```go
type Order struct {
    PlacedAt   models.Timestamp      // UTC → Bangkok on scan
    Category   models.EnumScan[string]
    Details    models.JsonScan[map[string]any]
}

client.GetClient().GetContext(ctx, &order, "SELECT * FROM orders WHERE id = ?", 1)
// Timestamp.Scan() auto-converts UTC→Bangkok; EnumScan and JsonScan map directly
```

---

## Edge Cases & Gotchas (Critical for Agents)

### 1. Timestamp UTC Conversion Is Unidirectional

```go
// Incoming: UTC timestamp from external API → stored as Bangkok time ✅
ts.UnmarshalJSON([]byte(`"2025-04-03T12:05:12Z"`))
ts.String() // "2025-04-03 19:05:12" (UTC+7h)

// Outgoing: Always displayed in Bangkok time ✅
json.Marshal(ts) // → "2025-04-03 19:05:12"

// Incoming: Local timezone (non-UTC) → stored as-is ✅
ts.Scan(time.Date(2024, 6, 15, 14, 30, 45, 0, time.Local))
// No conversion — uses the local value directly
```

### 2. Timestamp.ToUnix() Returns UTC Seconds

`ToUnix()` calls `time.Time.Unix()` internally which always returns UTC-based seconds. The stored `Timestamp` is in Bangkok location but Unix epoch is timezone-independent:

```go
ts := NewTimestampFromString("2025-04-03 19:05:12")  // This was "12:05:12Z" UTC
ts.ToUnix()                                            // Same as time.Unix(12:05:12Z) — always correct
```

### 3. EnumScan.MarshalJSON(nil) Returns Empty Byte Slice, Not Error

Unlike most Go types that error on nil JSON marshal, `EnumScan[T]` returns `[]byte("")`. This is intentional for nullable enum columns but means consumers should check `IsZero()` before expecting a value.

### 4. JsonScan.Data() Never Returns Nil — It Allocates

```go
var js JsonScan[map[string]interface{}]
data := js.Data()  // Allocates zero map, returns it — never panics
data["key"] = "value"  // Safe to use immediately
```

### 5. JsonScan with Struct Types — Strongly Typed vs Map

**`JsonScan[T]` works with any `T`** — not just raw maps:

| T Type | When to Use | Example |
|--------|-------------|---------|
| `map[string]interface{}` | Flexible/unknown schema, dynamic data | `{"key": "value", "nested": {...}}` |
| `[]interface{}` | Generic array without typed elements | `[1, 2, 3]`, `["a", "b"]` |
| Typed struct (`Profile`) | Known schema, compile-time safety | `Profile{Name: "Bob", Bio: "..."}` |
| Slice of pointers (`[]*Permission`) | Many-to-many relations stored as JSON | `[{"id": "uuid", "name": "read"}]` |

```go
type Permission struct {
    ID    uuid.UUID          `json:"id"`
    Name  string             `json:"name"`
}

type Profile struct {
    Bio      string            `json:"bio"`
    Settings map[string]string `json:"settings"`
}

// Option A: Map — flexible but no type safety
var userA User
userA.Settings = *NewJsonScan(map[string]interface{}{"theme": "dark"})
val := userA.Settings.Data().(map[string]interface{})  // type assertion required

// Option B: Struct — compile-time safety, direct field access
type User struct { Settings JsonScan[Profile] }
var userB User
userB.Settings = *NewJsonScan(Profile{Bio: "Hello", Settings: map[string]string{"theme": "dark"}})
bio := userB.Settings.Data().Bio  // direct access, no type assertion!

// Option C: Slice of pointers — many-to-many as JSON column
type User struct { Permissions JsonScan[[]*Permission] }
var userC User
userC.Permissions = *NewJsonScan([]*Permission{
    {ID: uuid.New(), Name: "read"},
    {ID: uuid.New(), Name: "write"},
})
perms := userC.Permissions.Data()  // returns []*Permission
for _, p := range perms {
    fmt.Println(p.Name)  // "read", "write"
}
```

### 6. Empty/Null JSON Unmarshal Behavior

| Type | Input `[]byte("null")` | Input `[]byte("")` | Input `[]byte(`""`)` |
|------|------------------------|-------------------|--------------------|
| **Date** | No change (keeps old value) | No change | No change |
| **Timestamp** | No change (keeps old value) | No change | No change |
| **EnumScan[T]** | Error: empty slice | No change | Unmarshal string to T |
| **JsonScan[T]** | No change | No change | Unmarshal content |

### 7. Panic vs Error Constructors

- `NewDateFromString()` — **panics** on invalid input
- `ParseDateFromString()` — returns error (safe)
- Same pattern for `Timestamp`: `NewTimestampFromString()` panics, `ParseTimestampFromString()` returns error

**Agent guidance**: Use `Parse*` in production code where errors should be handled; use `New*` in tests or initialization where invalid input is truly impossible.

### 8. Zero Value Behavior in SQL

- `Date{}.Value()` → returns `nil` (useful for nullable columns — writes NULL)
- `Timestamp{}.Value()` → returns `nil` (same pattern)
- `EnumScan[T]` with nil v → `Value()` returns `nil`
- `JsonScan[T]` with nil v → `Value()` returns `nil`

---

## Testing Reference

Test files are located alongside source files and cover all edge cases:

| Test File | Tests | Key Edge Cases Covered |
|-----------|-------|----------------------|
| `date_test.go` | Date constructors, JSON marshal/unmarshal, SQL scan/value, week/weekday | nil scan, empty string, invalid format, zero date value |
| `timestamp_test.go` | Timestamp constructors, UTC→Bangkok conversion, multiple input formats, SQL scan | all 14+ JSON formats, nil scan, type mismatch error, ValueOrZero |
| `enum_test.go` | EnumScan generics, Data/Set/String/IsZero, JSON marshal/unmarshal | nil pointer handling, zero string, empty slice on marshal |
| `json_test.go` | JsonScan for maps/arrays/structs, SQL scan with auto-detection, fluent Set() | nil slice marshals to `[]`, invalid JSON error, object vs array detection, typed structs (Profile, Permission) |

```bash
go test ./... -v
```

---

## Entity Naming Conventions (for Graphify / Knowledge Indexing)

When analyzing or referencing entities from this package:

- Use lowercase with underscores: `date_type`, `timestamp_utc_conversion`, `enumscan_null_handling`
- One level of parent only: `{filename}_{entity}` — e.g., `date_date` for the Date type in `date.go`
- Module-level entities use module prefix: `models_timestamp`, `models_enums_can`

---

## Agent Workflow Guidelines

### When Modifying a Type

1. **Read this AGENTS.md first** → understand the type's contract (JSON/SQL behavior, edge cases)
2. **Read the source file** (`date.go`, `timestamp.go`, etc.) → verify current implementation
3. **Read corresponding `_test.go`** → check which edge cases are already covered
4. **Update tests first** → if adding a new feature or changing behavior
5. **Update README.md** → keep API reference aligned with actual code
6. **Run tests**: `go test ./... -v`

### When Adding a New Type

1. Follow the same pattern: wrap underlying type, implement Scanner/Valuer/Marshaler/Unmarshaler
2. Use constructor naming convention: `New{Name}From*()` for panicking, `Parse{Name}From*()` for error-returning
3. Add tests covering: zero value, nil scan, invalid input, JSON marshal/unmarshal both directions
4. Document in README.md with same format as existing types

### When Debugging an Issue

1. **Wrong timezone?** → Check if the incoming value has `location == time.UTC` — Timestamp auto-converts, Date does not
2. **SQL scan failing?** → Verify the source type matches: `time.Time`, `string`, or `[]byte` only
3. **JSON marshal unexpected?** → Timestamp always outputs Bangkok time; nil EnumScan returns empty slice not error
4. **Panicking in production?** → Check if `New*()` constructor is used instead of `Parse*()`

---

## Related Modules Reference

| Module | Description | Usage Pattern |
|--------|-------------|---------------|
| **[fn](../fn/README.md)** | Generic struct-to-struct conversion via JSON round-trip | Use `ConvertStruct` or `CopyJSON` to map models types between database and API layers when field names differ |
| **psql** | PostgreSQL connection pool with tracing hooks | Types implement Scanner/Valuer for direct column mapping — no manual conversion needed |

---

*Last updated: 2026-06-14*
