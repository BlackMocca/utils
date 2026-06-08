package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ======================== Scan tests ========================

func TestEnumScan_Scan_String(t *testing.T) {
	var e EnumScan[string]
	err := e.Scan("active")
	assert.NoError(t, err)
	assert.Equal(t, "active", e.Data())
}

func TestEnumScan_Scan_Nil(t *testing.T) {
	var e EnumScan[string]
	err := e.Scan(nil)
	assert.NoError(t, err)
	assert.True(t, e.IsZero())
}

func TestEnumScan_Set_Get(t *testing.T) {
	var e EnumScan[string]
	e.Set("pending")
	assert.Equal(t, "pending", e.Data())
}

func TestEnumScan_NewEnumScan(t *testing.T) {
	e := NewEnumScan[string]("active")
	assert.Equal(t, "active", e.Data())
	assert.False(t, e.IsZero())
}

func TestEnumScan_Value_NonNil(t *testing.T) {
	e := NewEnumScan[string]("active")
	val, err := e.Value()
	assert.NoError(t, err)
	assert.Equal(t, "active", val)
}

func TestEnumScan_Value_Nil(t *testing.T) {
	var e EnumScan[string]
	val, err := e.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

// ======================== MarshalJSON tests ========================

func TestEnumScan_MarshalJSON_NonNil(t *testing.T) {
	e := NewEnumScan[string]("active")
	b, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, `"active"`, string(b))
}

func TestEnumScan_MarshalJSON_Nil(t *testing.T) {
	var e EnumScan[string]
	// When v is nil, MarshalJSON returns []byte("") which is invalid JSON,
	// causing json.Marshal to error with "unexpected end of JSON input"
	b, err := json.Marshal(e)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")
	assert.Equal(t, ``, string(b))
}

// ======================== UnmarshalJSON tests ========================

func TestEnumScan_UnmarshalJSON_Valid(t *testing.T) {
	var e EnumScan[string]
	input := []byte(`"active"`)
	err := json.Unmarshal(input, &e)
	assert.NoError(t, err)
	assert.Equal(t, "active", e.Data())
}

func TestEnumScan_UnmarshalJSON_Empty(t *testing.T) {
	var e EnumScan[string]
	input := []byte("")
	// Empty input causes json.Unmarshal to error with "unexpected end of JSON input"
	err := json.Unmarshal(input, &e)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")
}

func TestEnumScan_UnmarshalJSON_WithExistingValue(t *testing.T) {
	e := NewEnumScan[string]("pending")
	input := []byte(`"active"`)
	err := json.Unmarshal(input, &e)
	assert.NoError(t, err)
	assert.Equal(t, "active", e.Data())
}

// ======================== String / IsZero tests ========================

func TestEnumScan_String_NonNil(t *testing.T) {
	e := NewEnumScan[string]("active")
	assert.Equal(t, "active", e.String())
}

func TestEnumScan_String_Nil(t *testing.T) {
	var e EnumScan[string]
	assert.Equal(t, "", e.String())
}

func TestEnumScan_IsZero_True(t *testing.T) {
	var e EnumScan[string]
	assert.True(t, e.IsZero())
}

func TestEnumScan_IsZero_False(t *testing.T) {
	e := NewEnumScan[string]("active")
	assert.False(t, e.IsZero())
}

// ======================== Type parameter tests ========================

func TestEnumScan_String_Type(t *testing.T) {
	e := NewEnumScan[string]("admin")
	assert.Equal(t, "admin", e.Data())
	b, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, `"admin"`, string(b))
}
