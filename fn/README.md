# `fn` — Generic Struct & JSON Utility Package

[![Go Reference](https://pkg.go.dev/badge/github.com/BlackMocca/utils/fn.svg)](https://pkg.go.dev/github.com/BlackMocca/utils/fn)
[![Test](https://github.com/BlackMocca/utils/workflows/test/badge.svg)](https://github.com/BlackMocca/utils/actions)

Lightweight Go utilities for converting between structs via JSON serialization. Built with **Go 1.26+ generics**, no external runtime dependencies beyond the standard library.

---

## Install

```bash
go get github.com/BlackMocca/utils/fn
```

---

## Overview

This package provides two functions that marshal a value to JSON and unmarshal it into another type — effectively copying data between structs even when their field **names differ** (as long as the `json` tags align).

| Function | Signature | Return | Use Case |
|---|---|---|---|
| **`ConvertStruct`** | `func ConvertStruct[Src any, Dst any](src Src) (Dst, error)` | `(Dst, error)` | Generic API; destination type is inferred by the caller. No pointer needed on output. |
| **`CopyJSON`** | `func CopyJSON(src any, dst any) error` | `error` | Classic Go pattern — accepts a source value and a **pointer** to the destination. |

Both functions rely on `encoding/json` under the hood, so all standard JSON rules apply (tags, omitempty, slices, maps, nested structs, pointers, etc.).

---

## Examples

### ConvertStruct — generic, zero-cost destination

```go
type SourceUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DestUser struct {
	ID       int    `json:"id"`
	FullName string `json:"name"` // renamed field (JSON tag matches)
	Contact  string `json:"email"`
}

src := SourceUser{ID: 1, Name: "Alice", Email: "alice@example.com"}
dst, err := fn.ConvertStruct[SourceUser, DestUser](src)
if err != nil {
    // handle error
}
fmt.Printf("%+v\n", dst)
// Output: {ID:1 FullName:Alice Contact:alice@example.com}
```

### CopyJSON — pointer-based (classic pattern)

```go
var dst DestUser
err := fn.CopyJSON(src, &dst)
if err != nil {
    // handle error
}
```

### Nested structs, pointers, slices, maps — all supported

```go
type InnerMeta struct {
    Tags  []string `json:"tags"`
    Count int      `json:"count"`
}

type NestedSrc struct {
    ID       int       `json:"id"`
    Metadata InnerMeta `json:"metadata"`
}

type NestedDest struct {
    ID       int       `json:"id"`
    Metadata InnerMeta `json:"metadata"`
}

src := NestedSrc{
    ID: 7,
    Metadata: InnerMeta{
        Tags:  []string{"go", "test"},
        Count: 3,
    },
}

dst, err := fn.ConvertStruct[NestedSrc, NestedDest](src)
// dst now mirrors src with matching JSON-tagged fields preserved
```

### Omitted fields (`json:"-"`)

Fields tagged `json:"-"` are **excluded** during serialization and will hold the zero value in the destination.

```go
type BadSrc struct {
    Value interface{} `json:"-"` // omitted from JSON round-trip
}

src := BadSrc{Value: "ignored"}
dst, err := fn.ConvertStruct[BadSrc, BadSrc](src)
// dst.Value == nil
```

---

## How It Works

Both functions follow the same two-step process internally:

1. **`json.Marshal(src)`** — serializes the source value to a JSON byte slice.
2. **`json.Unmarshal(bytes, &dst)`** — deserializes that JSON into the destination type.

Because `ConvertStruct` uses Go 1.26+ generics (`Src any, Dst any`), the caller specifies both types at compile time and gets back a **non-pointer value**. `CopyJSON` mirrors the traditional pattern where the caller passes a pointer to the destination variable.

---

## When to Use (and When Not To)

### ✅ Good for
- Mapping between API request/response DTOs with different field names.
- Deep-copying structs that aren't safe for shallow copy.
- Transforming data layers (domain → view model, proto → struct, etc.).
- Any scenario where JSON round-trip is acceptable and desired.

### ❌ Not suitable for
- Performance-critical paths that cannot afford marshalling overhead.
- Structs with `time.Time` fields without custom `json.Unmarshaler` implementations (the standard library handles `time.Time` fine, but edge cases may arise).
- Circular references — will cause infinite recursion during marshal.

---

## Testing

```bash
go test ./... -v
```

The test suite covers:
- Same-type and cross-type conversion
- Different JSON tag names (renamed fields)
- Nested structs, slices, maps, and pointers (including `nil` pointers)
- Empty structs and zero-value inputs
- `json:"-"` omitted fields
- Round-trip consistency between both functions

---

## License

This project is licensed under the terms of the [LICENSE](../LICENSE) file in the parent repository.
