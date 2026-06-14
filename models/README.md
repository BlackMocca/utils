# `models` — Database & JSON Serialization Types with Automatic Timezone Handling

**Package**: `github.com/BlackMocca/utils/models`  
**Go Version**: 1.26+ (uses generics for `EnumScan`, `JsonScan`)  
**Dependencies**: `4d63.com/tz`

Custom types for database and JSON serialization with automatic **Asia/Bangkok (ICT, UTC+7)** timezone handling on `Timestamp`. All four types implement `database/sql.Scanner` and `driver.Valuer` for direct column mapping, plus `json.Marshaler`/`Unmarshaler` for seamless JSON wire format.

---

## Install

```bash
go get github.com/BlackMocca/utils/models
```

---

## Overview

This package provides four custom types that wrap Go's standard library types with additional serialization logic:

| Type | Wraps | Purpose | Key Feature |
|------|-------|---------|-------------|
| **Date** | `time.Time` | Date-only operations (`YYYY-MM-DD`) | No timezone conversion; format always `2006-01-02` |
| **Timestamp** | `time.Time` | Full timestamp with automatic UTC→Bangkok conversion | Auto-conversion on JSON unmarshal & SQL scan |
| **EnumScan[T]** | `*T` (string-like) | Nullable enum/string with DB support | Zero-value safe: `Data()` never panics, `IsZero()` checks both nil and empty string |
| **JsonScan[T]** | `*T` (any type) | Generic JSON column — maps, arrays, or structs | Auto-detects object vs array during SQL scan; fluent `Set()` builder pattern |

### ⚠️ Critical Timezone Rule

**Timestamp automatically converts UTC → Asia/Bangkok (UTC+7)** during both JSON unmarshaling and SQL scanning. Local-timezone values are stored as-is without conversion.

```
Input: "2025-04-03T12:05:12Z"  (UTC)
Stored & displayed: "2025-04-03 19:05:12"  (Bangkok, UTC+7)

Input: "2025-04-03 14:00:00"   (local/non-UTC)
Stored & displayed: "2025-04-03 14:00:00"  (no conversion)
```

---

## Date

Wraps `time.Time` for date-only operations (`YYYY-MM-DD`). No timezone conversion.

### Constructors

| Function | Description | Example |
|----------|-------------|---------|
| `NewDateFromString(date string)` | Create from `"2024-06-15"` — panics on invalid input | `d := NewDateFromString("2024-06-15")` |
| `ParseDateFromString(date string)` | Same but returns error — safe for production | `d, err := ParseDateFromString("2024-06-15")` |
| `NewDateFromTime(t time.Time)` | Create from standard time | `d := NewDateFromTime(time.Now())` |

### Methods

```go
d.String()                    // → "2024-06-15" (always YYYY-MM-DD)
d.Format("Mon Jan 2 2006")    // Custom format: "Sat Jun 15 2024"
d.ToTime() time.Time          // Convert back to standard time.Time
d.ToPointer() *Date           // Get pointer for nullable fields
d.ToTimestamp() Timestamp     // Promote date → full timestamp (time = 00:00:00 Bangkok)
d.Weekday() time.Weekday      // Returns day of week (Sunday=0, Monday=1, …)
```

### JSON Marshaling

```go
// Marshal — always outputs "YYYY-MM-DD"
b, _ := json.Marshal(d)  // → `"2024-06-15"`

// Unmarshal — auto-detects YYYY-MM-DD format; null/empty/"nil" = no change
var d Date
json.Unmarshal([]byte(`"2024-06-15"`), &d)    // Success
json.Unmarshal([]byte(`null`), &d)             // No change to existing value
json.Unmarshal([]byte(`""`), &d)               // No change (empty string ignored)
json.Unmarshal([]byte(`"invalid"`), &d)        // Error: invalid format
```

### SQL Scan / Value

```go
var d Date

// Scan from database — accepts time.Time, string, or []byte
err := d.Scan("2024-06-15")                     // string
err = d.Scan([]byte("2024-06-15"))              // []byte
err = d.Scan(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)) // time.Time

// Value for writing — returns "YYYY-MM-DD" or nil if zero date
val, err := d.Value()     // → `"2024-06-15"` or nil (zero date writes NULL)

// Edge cases:
d.Scan(nil)               // No error, keeps zero value
d.Scan("not-a-date")      // Error returned
```

---

## Timestamp

Wraps `time.Time` with full timestamp support and **automatic UTC→Bangkok conversion**.

### Constructors

| Function | Description | Example |
|----------|-------------|---------|
| `NewTimestampFromNow()` | Current time in Bangkok | `ts := NewTimestampFromNow()` |
| `NewTimestampFromString(s string)` | From `"2024-06-15 14:30:45"` — panics on invalid input | `ts := NewTimestampFromString("2024-06-15 14:30:45")` |
| `ParseTimestampFromString(s string)` | Same but returns error — safe for production | `ts, err := ParseTimestampFromString("2024-06-15 14:30:45")` |
| `NewTimestampFromTime(t time.Time)` | Create from standard time (converts to Bangkok) | `ts := NewTimestampFromTime(time.Now())` |

### Methods

```go
ts.String()                    // → "2024-06-15 14:30:45" (always Bangkok time)
ts.Format("2006/01/02")        // Custom format: "2024/06/15"
ts.ToTime() time.Time          // Convert back to standard time.Time
ts.ToPointer() *Timestamp      // Get pointer for nullable fields
ts.ToUnix() int64              // Unix epoch seconds (timezone-independent)
ts.YearDay() int               // Day of year (1–365/366)
ts.ValueOrZero() string        // Returns formatted string or "" if zero (safe for nil)
```

### Supported JSON Input Formats (UnmarshalJSON)

Multiple timestamp formats are auto-detected during unmarshaling:

```go
"2025-04-03 12:05:12.510131"     // Microseconds
"2025-04-03 12:05:12.510"        // Milliseconds
"2025-04-03 12:05:12"            // Standard (space separator)
"2024-06-25T14:30:00"            // ISO8601 (no timezone)
"2024-06-25"                     // Date only
"time.RFC3339"                   // "2025-04-03T12:05:12Z"
"time.RFC3339Nano"               // With nanoseconds and Z suffix
"time.RFC1123"                   // "Mon, 25 Jun 2024 14:30:00 UTC"
"time.UnixDate", "time.ANSIC"    // Standard Go formats
```

### JSON Marshaling

```go
// Marshal — always outputs "YYYY-MM-DD HH:MM:SS" in Bangkok time
b, _ := json.Marshal(ts)  // → `"2025-04-03 19:05:12"` (if originally UTC noon)

// Unmarshal — auto-detects format, converts UTC → Bangkok
var ts Timestamp
json.Unmarshal([]byte(`"2025-04-03T12:05:12Z"`), &ts)  // UTC 12:05 → +7h = 19:05 Bangkok

// Null handling — null, "", "nil" all leave the field unchanged
```

### SQL Scan / Value

```go
var ts Timestamp

// Scan from database — accepts time.Time, string, or []byte
err := ts.Scan("2024-06-15 14:30:45")                                     // string
err = ts.Scan([]byte("2024-06-15 14:30:45"))                             // []byte

// UTC → Bangkok auto-conversion on scan from time.Time:
err = ts.Scan(time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC))         // → stored as 21:30:45 Bangkok

// Value for writing — returns formatted string or nil if zero
val, err := ts.Value()     // → `"2024-06-15 14:30:45"` or nil (zero writes NULL)

// Edge cases:
ts.Scan(nil)               // No error
ts.Scan(42)                // Error: "cannot scan type int into Timestamp"
```

---

## EnumScan[T] — Nullable Generic Enum

Generic wrapper for enum/string types with SQL and JSON support. The `T` constraint is `~string`, so it works with both plain `string` and typed aliases like `type Status string`.

### Usage

```go
// Create from string value
e := NewEnumScan[string]("active")

// Or use a custom type alias
type Role string
const (
    RoleAdmin  Role = "admin"
    RoleUser   Role = "user"
)
r := NewEnumScan[Role](string(RoleAdmin))
```

### Methods

| Method | Return | Description |
|--------|--------|-------------|
| `e.Data()` | `T` | Returns current value or zero type if nil — **never panics** |
| `e.Set(v T)` | `void` | Sets a new value |
| `e.String()` | `string` | `""` if nil, else the string value |
| `e.IsZero()` | `bool` | `true` if v is nil **or** empty string |

### JSON Marshaling

```go
b, _ := json.Marshal(e)  // → `"admin"` (non-nil) or `""` (nil — empty slice, not error!)

var e EnumScan[string]
json.Unmarshal([]byte(`"active"`), &e)      // Success — sets value to "active"
json.Unmarshal([]byte(""), &e)              // No change — empty input ignored
```

### SQL Scan / Value

```go
err := e.Scan("admin")     // string → success
err = e.Scan(nil)          // nil → zero value (no error), column becomes NULL

val, err := e.Value()       // Returns `"admin"` or nil if zero
```

---

## JsonScan[T] — Generic JSON Field

Generic wrapper for JSON data (maps, arrays, or structs) with SQL and JSON support. Auto-detects whether the JSON is an object or array during SQL scan.

### Usage

```go
// From a map (flexible, no type safety)
js := NewJsonScan(map[string]interface{}{"key": "value"})

// From a slice of primitives
js2 := NewJsonScan([]interface{}{1, 2, 3})

// From a typed struct (compile-time safety)
type Config struct { Color string }
js3 := NewJsonScan(Config{Color: "red"})

// From a slice of pointers — many-to-many JSON column
type Permission struct { ID uuid.UUID; Name string }
js4 := NewJsonScan([]*Permission{
    {ID: uuid.New(), Name: "read"},
    {ID: uuid.New(), Name: "write"},
})

// Fluent Set() — returns *JsonScan for chaining
js.Set(map[string]interface{}{"key": "new_value"}).Set(map[string]interface{}{"other": 123})
```

### Methods

| Method | Return | Description |
|--------|--------|-------------|
| `js.Data()` | `T` | Returns current value or allocates zero if nil — **never panics** |
| `js.Set(v T)` | `*JsonScan[T]` | Sets a new value, returns self for chaining |

### JSON Marshaling

```go
b, _ := json.Marshal(js)  // → `{"key":"value"}` or `[1,2,3]`

// Special: nil slices marshal to [] (empty array), not null
var js JsonScan[[]string]
json.Marshal(js)  // → `[]` — never `null`

var js2 JsonScan[map[string]interface{}]
js2.UnmarshalJSON([]byte(`{"name":"test"}`))  // Success

// Empty input is ignored (no error, no change)
```

### SQL Scan / Value

```go
var js JsonScan[map[string]interface{}]

err := js.Scan([]byte(`{"key":"value"}`))     // Object — auto-detected
err = js.Scan([]byte(`[1,2,3]`))              // Array — auto-detected
err = js.Scan(nil)                            // No error, no change

val, err := js.Value()    // Returns JSON-encoded []byte or nil if zero
```

---

## Common Patterns & Examples

### Pattern 1: Full Domain Model with All Four Types

```go
type User struct {
    ID        uint               
    CreatedAt models.Timestamp              // Auto UTC→Bangkok conversion on scan
    BirthDate models.Date                   // Date only, no timezone issues  
    Role      models.EnumScan[constants.Role]  // Nullable enum from DB
    Profile   models.JsonScan[map[string]interface{}]  // Flexible JSON column
}

// Usage with GORM or sqlx:
db.First(&user, id)
user.CreatedAt.String()  // "2024-06-15 14:30:45" (Bangkok time)
if user.Role.IsZero() { 
    fmt.Println("No role assigned") 
} else { 
    fmt.Printf("Role: %s\n", user.Role.Data()) 
}
```

### Pattern 2: Nullable Fields with Pointers

```go
type Article struct {
    PublishedAt *models.Timestamp  // nil = not yet published
    EffectiveDate *models.Date     // nil = no date restriction
}

// Safe zero access on nullable fields:
var article Article
fmt.Println(article.PublishedAt.ValueOrZero())  // "" (nil pointer, safe)
```

### Pattern 3: DTO Mapping Between Layers

Use generic struct-to-struct conversion functions to map between database models and API response DTOs when field names differ:

```go
// Database model (uses custom types)
type OrderDB struct {
    ID       int            `json:"id"`
    PlacedAt models.Timestamp `json:"placed_at"`
}

// API response (flat string fields)
type OrderAPI struct {
    OrderID int    `json:"order_id"`
    Time    string `json:"time"`
}

// Map via ConvertStruct — Timestamp.MarshalJSON() → "2024-06-15 14:30:45"
dst, err := fn.ConvertStruct[OrderDB, OrderAPI](dbModel)
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

js := JsonScan[map[string]interface{}]{}
data := js.Data()  // allocates zero map — safe to use immediately
data["key"] = "value"
```

---

## Edge Cases & Gotchas

1. **Timestamp UTC Conversion**  
   Incoming values with `location == time.UTC` are automatically converted to Bangkok (UTC+7). Local timezone values are stored as-is. Marshal and Value always output Bangkok time.

2. **Timestamp.ToUnix() Is Always Correct**  
   `ToUnix()` calls `time.Time.Unix()` internally, which returns seconds since epoch regardless of the time.Location — always the correct Unix timestamp.

3. **EnumScan Nil Marshal Returns Empty Slice**  
   `MarshalJSON(nil)` returns `[]byte("")` (empty slice), not an error. This is intentional for nullable enum columns but means consumers should check `IsZero()` before expecting a value.

4. **JsonScan Nil Slice Marshals to `[]`, Not `null`**  
   Nil slices marshal to empty array `[]` — this avoids the "null vs empty" ambiguity in JSON responses where an empty collection is semantically different from "no data".

5. **JsonScan with Struct Types — Strongly Typed vs Map**

   `JsonScan[T]` works with **any** `T` type:

   | T Type | When to Use | Example |
   |--------|-------------|---------|
   | `map[string]interface{}` | Flexible/unknown schema, dynamic data | `{"key": "value", "nested": {...}}` |
   | `[]interface{}` | Generic array without typed elements | `[1, 2, 3]`, `["a", "b"]` |
   | Typed struct (`Profile`) | Known schema, compile-time safety | `Profile{Name: "Bob", Bio: "..."}` |
   | Slice of pointers (`[]*Permission`) | Many-to-many relations stored as JSON | `[{"id": "uuid", "name": "read"}]` |

   ```go
   // Map — flexible but no type safety (requires type assertion)
   var userA User
   val := userA.Settings.Data().(map[string]interface{})

   // Struct — compile-time safety, direct field access
   bio := userB.Profile.Data().Bio  // no type assertion needed!

   // Slice of pointers — many-to-many as JSON column
   perms := userC.Permissions.Data()  // returns []*Permission
   ```

6. **Date Zero Value Writes NULL**  
   A zero `Date{}` returns `nil` from `.Value()`, making it ideal for nullable DB columns. Same pattern applies to `Timestamp`.

7. **Panic vs Error Constructors**  
   `New*()` constructors panic on invalid input — use in tests or initialization. `Parse*()` functions return error — use in production code where input may be untrusted.

8. **String() vs Format()**
   - `String()` always uses the default layout (`DateLayout` or `TimestampLayout`)
   - `Format(layout)` accepts any Go time format string

---

## Testing

Test files are located alongside source files and cover all edge cases:

```bash
go test ./... -v
```

| Test File | Tests | Key Edge Cases Covered |
|-----------|-------|----------------------|
| `date_test.go` | Constructors, JSON marshal/unmarshal, SQL scan/value, weekday | nil scan, empty string, invalid format, zero date value |
| `timestamp_test.go` | Constructors, UTC→Bangkok conversion, 14+ input formats, SQL scan | all formats, nil scan, type mismatch error, ValueOrZero |
| `enum_test.go` | Generics, Data/Set/String/IsZero, JSON marshal/unmarshal | nil pointer handling, zero string, empty slice on marshal |
| `json_test.go` | Maps/arrays/structs, SQL auto-detection, fluent Set() | nil slice marshals to `[]`, invalid JSON error, object vs array detection, typed structs (Profile, Permission) |

---

## Related Modules Reference

- **[fn](../fn/README.md)** — Use `ConvertStruct` or `CopyJSON` to map models types between database and API layers when field names differ.
- **psql** — PostgreSQL client that works seamlessly with all four types via Scanner/Valuer interfaces. All custom types implement `database/sql.Scanner` and `driver.Valuer` for direct column mapping without manual conversion.

---

*Last updated: 2026-06-14*
