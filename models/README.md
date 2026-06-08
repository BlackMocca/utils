# models — Database & JSON Serialization Types

Install with:

```bash
go get github.com/BlackMocca/utils/models
```

**Package**: `github.com/BlackMocca/utils/models`  
**Go Version**: 1.26+  
**Dependencies**: `4d63.com/tz` (timezone support)

## Overview

This package provides custom types for database and JSON serialization with automatic **Asia/Bangkok timezone handling**. All date/time types store data in ICT (UTC+7).

---

## ⚠️ Timezone Behavior

- Default location: `Asia/Bangkok` (ICT, UTC+7)
- **Timestamp** automatically converts UTC → Bangkok during parsing/scan
- Date types do not perform timezone conversion (stored as-is)

---

## 1. Date Type

Wraps `time.Time` for date-only operations (`YYYY-MM-DD`).

```go
type Date time.Time
```

### Constructors

| Function | Description | Example |
|----------|-------------|---------|
| `NewDateFromString(date string)` | Create from `"2024-06-15"` | `d := NewDateFromString("2024-06-15")` |
| `ParseDateFromString(date string)` | Same but returns error | `d, err := ParseDateFromString("2024-06-15")` |
| `NewDateFromTime(t time.Time)` | Create from standard time | `d := NewDateFromTime(time.Now())` |

### Key Methods

```go
d.String()                    // Returns "2024-06-15"
d.Format("Mon Jan 2 2006")    // Custom format: "Sat Jun 15 2024"
d.ToTime() time.Time          // Convert to standard time.Time
d.ToPointer() *Date           // Get pointer to Date
d.ToTimestamp() Timestamp     // Convert to Timestamp type
d.Weekday() time.Weekday      // Returns day of week (Sunday=0)
```

### JSON Marshaling

```go
// Marshal
b, _ := json.Marshal(d)  // → "2025-06-15"

// Unmarshal
var d Date
json.Unmarshal([]byte(`"2024-06-15"`), &d)    // Success
json.Unmarshal([]byte(`null`), &d)             // Null → no change
json.Unmarshal([]byte(`""`), &d)               // Empty → no change
json.Unmarshal([]byte(`"invalid"`), &d)        // Returns error

// Unmarshal handles these special values without error:
// null, "", "nil"
```

### Scan / Value (SQL)

```go
var d Date

// Scan from database
err := d.Scan("2024-06-15")                     // string
err = d.Scan([]byte("2024-06-15"))              // []byte
err = d.Scan(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)) // time.Time

// Value to database
val, err := d.Value()     // Returns "2024-06-15" or nil if zero date

// Edge cases:
d.Scan(nil)               // No error, keeps zero value
d.Scan("not-a-date")      // Error returned
```

---

## 2. Timestamp Type

Wraps `time.Time` with full timestamp support and **UTC→Bangkok conversion**.

```go
type Timestamp time.Time
```

### Constructors

| Function | Description | Example |
|----------|-------------|---------|
| `NewTimestampFromNow()` | Current time | `ts := NewTimestampFromNow()` |
| `NewTimestampFromString(s)` | From `"2024-06-15 14:30:45"` | `ts := NewTimestampFromString("2024-06-15 14:30:45")` |
| `ParseTimestampFromString(s)` | Same but returns error | `ts, err := ParseTimestampFromString("...")` |

### Supported Input Formats (UnmarshalJSON)

The following formats are auto-detected during JSON unmarshaling:

```go
"2025-04-03 12:05:12.510131"     // Microseconds
"2025-04-03 12:05:12.510"        // Milliseconds
"2025-04-03 12:05:12"            // Standard
"2024-06-25T14:30:00"            // ISO8601 (no timezone)
"2024-06-25"                     // Date only
"time.RFC3339"                   // "2025-04-03T12:05:12Z"
"time.RFC3339Nano"               // With nanoseconds
"time.RFC1123"                   // "Mon, 25 Jun 2024 14:30:00 UTC"
"time.UnixDate", "time.ANSIC"    // Standard Go formats
```

### Key Methods

```go
ts.String()                    // Returns "2024-06-15 14:30:45"
ts.Format("2006/01/02")        // Custom format: "2024/06/15"
ts.ToTime() time.Time          // Convert to standard time.Time
ts.ToPointer() *Timestamp      // Get pointer
ts.ToUnix() int64              // Unix timestamp (seconds)
ts.YearDay() int               // Day of year (1-365/366)
ts.ValueOrZero() string        // Returns formatted string or "" if zero
```

### JSON Marshaling

```go
// Marshal - always outputs "2006-01-02 15:04:05" format (Bangkok time)
b, _ := json.Marshal(ts)  // → "2025-04-03 19:05:12"

// Unmarshal - auto-detects format, converts UTC → Bangkok
var ts Timestamp
json.Unmarshal([]byte(`"2025-04-03T12:05:12Z"`), &ts)  // UTC → +7h = 19:05:12

// Null handling (same as Date): null, "", "nil" → no change
```

### Scan / Value (SQL)

```go
var ts Timestamp

// Scan from database
err := ts.Scan("2024-06-15 14:30:45")                                     // string
err = ts.Scan([]byte("2024-06-15 14:30:45"))                             // []byte
err = ts.Scan(time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC))          // UTC → Bangkok auto-conversion
err = ts.Scan(time.Date(2024, 6, 15, 14, 30, 45, 0, Location))          // Local → no conversion

// Value to database
val, err := ts.Value()     // Returns "2024-06-15 14:30:45" or nil if zero

// Edge cases:
ts.Scan(nil)               // No error
ts.Scan(42)                // Error: "cannot scan type int into Timestamp"
```

---

## 3. EnumScan Generic Type

Generic wrapper for enum/string types with nullable support.

```go
type EnumScan[T ~string] struct { /* internal */ }
```

### Usage

```go
// Create from string value
e := NewEnumScan[string]("active")
e := NewEnumScan[StatusType]("pending")  // Custom type alias

// Get/Set data
val := e.Data()         // Returns current value or zero
e.Set("new_value")      // Set new value

// String representation
str := e.String()       // "" if nil, else the value
bool := e.IsZero()      // true if v is nil or empty string
```

### JSON Marshaling

```go
b, err := json.Marshal(e)  // → "active" (non-nil) or error (nil)

var e EnumScan[string]
json.Unmarshal([]byte(`"active"`), &e)      // Success
json.Unmarshal([]byte(""), &e)              // Error: empty input
```

### Scan / Value (SQL)

```go
err := e.Scan("active")     // string → success
err = e.Scan(nil)           // nil → zero value, no error

val, err := e.Value()       // Returns "active" or nil if zero
```

---

## 4. JsonScan Generic Type

Generic wrapper for JSON data (maps/arrays/structs) with database support.

```go
type JsonScan[T any] struct { /* internal */ }
```

### Usage

```go
// Create from value
js := NewJsonScan(map[string]interface{}{"key": "value"})
js := NewJsonScan([]interface{}{1, 2, 3})
js := NewJsonScan(MyStruct{Name: "test"})

// Get/Set data
data := js.Data()           // Returns T (allocates zero if nil)
js.Set(newValue)            // Sets new value (returns *JsonScan for chaining)
```

### JSON Marshaling

```go
b, _ := json.Marshal(js)  // → {"key":"value"}

var js JsonScan[map[string]interface{}]
json.Unmarshal([]byte(`{"name":"test"}`), &js)    // Success
json.Unmarshal([]byte(""), &js)                   // Error: empty input

// Special: nil slices marshal to [] not null
```

### Scan / Value (SQL)

```go
var js JsonScan[map[string]interface{}]
err := js.Scan([]byte(`{"key":"value"}`))     // Object
err = js.Scan([]byte(`[1,2,3]`))              // Array auto-detected
err = js.Scan(nil)                            // No error
err = js.Scan([]byte("invalid"))              // Error

val, err := js.Value()    // Returns []byte (JSON-encoded) or nil if zero
```

---

## Common Patterns & Examples

### Pattern 1: Database Model with Custom Types

```go
type User struct {
    ID        uint               `gorm:"primary_key"`
    CreatedAt Timestamp          `gorm:"column:created_at;type:timestamp"`
    BirthDate Date                `gorm:"column:birth_date;type:date"`
    Status    EnumScan[StatusType]
    Metadata  JsonScan[map[string]interface{}]
}

// Usage with GORM
db.First(&user, id)
user.CreatedAt.String()  // "2024-06-15 14:30:45"
```

### Pattern 2: JSON API Response

```go
type ApiResponse struct {
    Date   Date       `json:"date"`
    Time   Timestamp  `json:"timestamp"`
}

resp := ApiResponse{
    Date:   NewDateFromString("2024-06-15"),
    Time:   NewTimestampFromNow(),
}
b, _ := json.Marshal(resp)
// → {"date":"2024-06-15","timestamp":"2024-06-15 14:30:45"}
```

### Pattern 3: Nullable Enum in Database

```go
type Role struct {
    Name EnumScan[string] `gorm:"column:name"`
}

var r Role
db.First(&r, "id=1")   // Scans from DB

if r.Name.IsZero() {
    fmt.Println("No role assigned")
} else {
    fmt.Printf("Role: %s\n", r.Name.Data())
}
```

### Pattern 4: JSON Field Handling

```go
type Product struct {
    ID     int
    Config JsonScan[map[string]interface{}] `gorm:"column:config;type:jsonb"`
}

// Insert
p := NewJsonScan(map[string]interface{}{
    "color": "red",
    "size":  "L",
})
product.Config = *p
db.Create(&product)

// Query & Update
var product Product
db.First(&product, id)
config := product.Config.Data()
config["color"] = "blue"
product.Config.Set(config)
db.Save(&product)
```

---

## Edge Cases & Gotchas

1. **Timestamp UTC Conversion**  
   When scanning/unmarshaling UTC timestamps, they're automatically converted to Bangkok (UTC+7). Local timezone values are not converted.

2. **EnumScan Nil Handling**  
   - `MarshalJSON(nil)` returns error (empty byte slice is invalid JSON)
   - `UnmarshalJSON("")` returns error

3. **JsonScan Slice Marshal**  
   nil slices marshal to `[]` (empty array), not `null`

4. **Date Zero Values**  
   Zero date (`0001-01-01`) returns `nil` from `.Value()`, useful for nullable DB columns

5. **String() vs Format()**  
   - `String()` always uses default layout
   - `Format(layout)` accepts any Go time format string

---

## Testing Reference

Test files demonstrate all edge cases:

```bash
# Run all tests
go test ./models/... -v

# Run specific type tests
go test ./models/... -run TestTimestamp -v
```
