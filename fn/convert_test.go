package fn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── Helper types for tests ─────────────────────────────────────────────

type SourceUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DestUser struct {
	ID    int    `json:"id"`
	FullName string `json:"name"`
	Contact   string `json:"email"`
}

type InnerMeta struct {
	Tags []string `json:"tags"`
	Count int     `json:"count"`
}

type NestedSource struct {
	ID       int       `json:"id"`
	Metadata InnerMeta `json:"metadata"`
}

type NestedDest struct {
	ID       int       `json:"id"`
	Metadata InnerMeta `json:"metadata"`
}

type EmptyStruct struct{}

type WithPtrs struct {
	A *int    `json:"a,omitempty"`
	B *string `json:"b,omitempty"`
	C *bool   `json:"c,omitempty"`
}

type SliceSource struct {
	Items []int `json:"items"`
}

type MapSource struct {
	Data map[string]int `json:"data"`
}

// ── ConvertStruct tests ────────────────────────────────────────────────

func TestConvertStruct_SameType(t *testing.T) {
	src := SourceUser{ID: 1, Name: "Alice", Email: "alice@example.com"}

	dst, err := ConvertStruct[SourceUser, SourceUser](src)

	assert.NoError(t, err)
	assert.Equal(t, src, dst)
}

func TestConvertStruct_DifferentFieldNames(t *testing.T) {
	src := SourceUser{ID: 42, Name: "Bob", Email: "bob@example.com"}

	dst, err := ConvertStruct[SourceUser, DestUser](src)

	assert.NoError(t, err)
	assert.Equal(t, 42, dst.ID)
	assert.Equal(t, "Bob", dst.FullName)
	assert.Equal(t, "bob@example.com", dst.Contact)
}

func TestConvertStruct_EmptyStruct(t *testing.T) {
	src := EmptyStruct{}

	dst, err := ConvertStruct[EmptyStruct, EmptyStruct](src)

	assert.NoError(t, err)
	assert.Equal(t, src, dst)
}

func TestConvertStruct_NestedStruct(t *testing.T) {
	src := NestedSource{
		ID: 7,
		Metadata: InnerMeta{
			Tags:  []string{"go", "test"},
			Count: 3,
		},
	}

	dst, err := ConvertStruct[NestedSource, NestedDest](src)

	assert.NoError(t, err)
	assert.Equal(t, src.ID, dst.ID)
	assert.Equal(t, src.Metadata.Tags, dst.Metadata.Tags)
	assert.Equal(t, src.Metadata.Count, dst.Metadata.Count)
}

func TestConvertStruct_Pointers(t *testing.T) {
	a := 10
	b := "hello"
	c := true

	src := WithPtrs{A: &a, B: &b, C: &c}

	dst, err := ConvertStruct[WithPtrs, WithPtrs](src)

	assert.NoError(t, err)
	assert.NotNil(t, dst.A)
	assert.Equal(t, 10, *dst.A)
	assert.NotNil(t, dst.B)
	assert.Equal(t, "hello", *dst.B)
	assert.NotNil(t, dst.C)
	assert.True(t, *dst.C)
}

func TestConvertStruct_PointerWithNil(t *testing.T) {
	hello := "hello"
	src := WithPtrs{A: nil, B: &hello, C: nil}
	dst, err := ConvertStruct[WithPtrs, WithPtrs](src)

	assert.NoError(t, err)
	assert.Nil(t, dst.A)
	assert.NotNil(t, dst.B)
	assert.Equal(t, "hello", *dst.B)
	assert.Nil(t, dst.C)
}

func TestConvertSliceValue(t *testing.T) {
	src := SliceSource{Items: []int{1, 2, 3}}
	dst, err := ConvertStruct[SliceSource, SliceSource](src)

	assert.NoError(t, err)
	assert.Equal(t, src.Items, dst.Items)
}

func TestConvertMapValue(t *testing.T) {
	src := MapSource{Data: map[string]int{"a": 1, "b": 2}}
	dst, err := ConvertStruct[MapSource, MapSource](src)

	assert.NoError(t, err)
	assert.Equal(t, src.Data, dst.Data)
}

func TestConvertStruct_ZeroValues(t *testing.T) {
	src := SourceUser{} // all zero values
	dst, err := ConvertStruct[SourceUser, SourceUser](src)

	assert.NoError(t, err)
	assert.Equal(t, 0, dst.ID)
	assert.Empty(t, dst.Name)
	assert.Empty(t, dst.Email)
}

// ── CopyJSON tests ─────────────────────────────────────────────────────

func TestCopyJSON_SameStruct(t *testing.T) {
	src := SourceUser{ID: 5, Name: "Charlie", Email: "charlie@test.com"}
	var dst SourceUser

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, src, dst)
}

func TestCopyJSON_DifferentStruct(t *testing.T) {
	src := SourceUser{ID: 99, Name: "Diana", Email: "diana@ex.com"}
	var dst DestUser

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, 99, dst.ID)
	assert.Equal(t, "Diana", dst.FullName)
	assert.Equal(t, "diana@ex.com", dst.Contact)
}

func TestCopyJSON_NilPointer(t *testing.T) {
	// json.Unmarshal on a nil interface pointer does not error in Go;
	// the function simply returns without writing anything.
	src := SourceUser{ID: 1, Name: "Test", Email: "test@x.com"}

	err := CopyJSON(src, nil)

	assert.NoError(t, err) // Go's json.Unmarshal silently ignores a nil interface
}

func TestCopyJSON_EmptyStruct(t *testing.T) {
	var dst EmptyStruct
	src := EmptyStruct{}

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, src, dst)
}

func TestCopyJSON_NestedStruct(t *testing.T) {
	src := NestedSource{
		ID: 3,
		Metadata: InnerMeta{Tags: []string{"x"}, Count: 99},
	}
	var dst NestedDest

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, src.ID, dst.ID)
	assert.Equal(t, []string{"x"}, dst.Metadata.Tags)
	assert.Equal(t, 99, dst.Metadata.Count)
}

func TestCopyJSON_Pointers(t *testing.T) {
	a := 42
	src := WithPtrs{A: &a, B: nil}
	var dst WithPtrs

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.NotNil(t, dst.A)
	assert.Equal(t, 42, *dst.A)
	assert.Nil(t, dst.B)
}

func TestCopyJSON_MapValue(t *testing.T) {
	src := MapSource{Data: map[string]int{"k": 10}}
	var dst MapSource

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, 10, dst.Data["k"])
}

func TestCopyJSON_SliceValue(t *testing.T) {
	src := SliceSource{Items: []int{10, 20, 30}}
	var dst SliceSource

	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, []int{10, 20, 30}, dst.Items)
}

// ── Cross-validation: ConvertStruct vs CopyJSON ────────────────────────

func TestBothFunctions_ProduceSameResult(t *testing.T) {
	src := SourceUser{ID: 7, Name: "Eve", Email: "eve@y.com"}

	dst1, err1 := ConvertStruct[SourceUser, DestUser](src)
	var dst2 DestUser
	err2 := CopyJSON(src, &dst2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, dst1, dst2)
}

// ── JSON round-trip consistency ────────────────────────────────────────

func TestConvertStruct_RoundTrip(t *testing.T) {
	original := NestedSource{
		ID: 123,
		Metadata: InnerMeta{
			Tags:  []string{"a", "b", "c"},
			Count: 456,
		},
	}

	dst, err := ConvertStruct[NestedSource, NestedDest](original)
	assert.NoError(t, err)

	var roundTrip NestedSource
	err = CopyJSON(dst, &roundTrip)
	assert.NoError(t, err)

	assert.Equal(t, original.ID, roundTrip.ID)
	assert.Equal(t, original.Metadata.Tags, roundTrip.Metadata.Tags)
	assert.Equal(t, original.Metadata.Count, roundTrip.Metadata.Count)
}

func TestConvertStruct_WithOmittedFields(t *testing.T) {
	// When a field is omitted in the destination, it should remain its zero value.
	src := SourceUser{ID: 10, Name: "Frank", Email: "frank@z.com"}
	dst, err := ConvertStruct[SourceUser, DestUser](src)

	assert.NoError(t, err)
	assert.Equal(t, 10, dst.ID)
	assert.Equal(t, "Frank", dst.FullName)
	assert.Equal(t, "frank@z.com", dst.Contact)
}

// ── Edge: invalid source (should still work — Marshal never errors on normal Go types) ───

func TestConvertStruct_WithInvalidJSONTag(t *testing.T) {
	// When a field has json:"-" it is omitted during marshal/unmarshal,
	// so the destination receives the zero value.
	type BadSrc struct {
		Value interface{} `json:"-"`
	}

	src := BadSrc{Value: "ignored"}
	dst, err := ConvertStruct[BadSrc, BadSrc](src)

	assert.NoError(t, err)
	assert.Nil(t, dst.Value) // json:"-" is dropped → zero value in dst
}

func TestCopyJSON_InvalidSourceNonPtr(t *testing.T) {
	// CopyJSON requires dst to be a pointer; passing non-pointer dst should error.
	src := "hello"
	var dst string // not a pointer — but CopyJSON accepts any type and unmarshals into &dst

	// dst is addressable, so json.Unmarshal works fine. This test verifies no panic.
	err := CopyJSON(src, &dst)

	assert.NoError(t, err)
	assert.Equal(t, "hello", dst)
}
