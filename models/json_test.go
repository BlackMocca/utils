package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ======================== Scan tests ========================

func TestJsonScan_Scan_ValidJSON_Object(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	input := []byte(`{"key":"value","number":42}`)
	err := js.Scan(input)
	assert.NoError(t, err)
	data := js.Data()
	assert.Equal(t, "value", data["key"])
	assert.Equal(t, float64(42), data["number"])
}

func TestJsonScan_Scan_ValidJSON_Array(t *testing.T) {
	var js JsonScan[[]interface{}]
	input := []byte(`[1,2,3]`)
	err := js.Scan(input)
	assert.NoError(t, err)
	data := js.Data()
	assert.Len(t, data, 3)
	assert.Equal(t, float64(1), data[0])
}

func TestJsonScan_Scan_Nil(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	err := js.Scan(nil)
	assert.NoError(t, err)
}

func TestJsonScan_Scan_EmptyBytes(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	input := []byte("")
	err := js.Scan(input)
	assert.NoError(t, err)
}

func TestJsonScan_Scan_InvalidJSON(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	input := []byte(`{invalid}`)
	err := js.Scan(input)
	assert.Error(t, err)
}

func TestJsonScan_Set_Get(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	js := NewJsonScan(data)
	result := js.Set(data)
	// Set returns *JsonScan[T], compare via Data()
	assert.Equal(t, data, result.Data())
}

// ======================== MarshalJSON tests ========================

func TestJsonScan_MarshalJSON_Map(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	js := NewJsonScan(data)
	b, err := json.Marshal(js)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(b, &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestJsonScan_MarshalJSON_Slice(t *testing.T) {
	data := []interface{}{1, 2, 3}
	js := NewJsonScan(data)
	b, err := json.Marshal(js)
	assert.NoError(t, err)
	assert.Equal(t, `[1,2,3]`, string(b))
}

func TestJsonScan_MarshalJSON_NilSlice(t *testing.T) {
	var data []interface{}
	js := NewJsonScan(data)
	b, err := json.Marshal(js)
	assert.NoError(t, err)
	// MarshalJSON checks nil slice but the check doesn't match, falls through to json.Marshal(nil) = null
	assert.Equal(t, `null`, string(b))
}

// ======================== UnmarshalJSON tests ========================

func TestJsonScan_UnmarshalJSON_Object(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	input := []byte(`{"name":"test","age":25}`)
	err := json.Unmarshal(input, &js)
	assert.NoError(t, err)
	data := js.Data()
	assert.Equal(t, "test", data["name"])
}

func TestJsonScan_UnmarshalJSON_Array(t *testing.T) {
	var js JsonScan[[]interface{}]
	input := []byte(`[10,20,30]`)
	err := json.Unmarshal(input, &js)
	assert.NoError(t, err)
	data := js.Data()
	assert.Len(t, data, 3)
}

func TestJsonScan_UnmarshalJSON_EmptyBytes(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	input := []byte("")
	// Empty input causes json.Unmarshal to error with "unexpected end of JSON input"
	err := json.Unmarshal(input, &js)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")
}

func TestJsonScan_UnmarshalJSON_ExistingData(t *testing.T) {
	existing := map[string]interface{}{"old": "value"}
	js := NewJsonScan(existing)
	input := []byte(`{"new":"data"}`)
	err := json.Unmarshal(input, &js)
	assert.NoError(t, err)
	data := js.Data()
	assert.Equal(t, "data", data["new"])
}

// ======================== Value tests ========================

func TestJsonScan_Value_NonNil(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	js := NewJsonScan(data)
	val, err := js.Value()
	assert.NoError(t, err)
	// Value returns []byte from json.Marshal
	expected, _ := json.Marshal(data)
	assert.Equal(t, expected, val)
}

func TestJsonScan_Value_Nil(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	val, err := js.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

// ======================== Round-trip tests ========================

func TestJsonScan_RoundTrip_MarshalJSON_UnmarshalJSON(t *testing.T) {
	data := map[string]interface{}{"key": "value", "nested": map[string]interface{}{}}
	js := NewJsonScan(data)

	b, err := json.Marshal(js)
	assert.NoError(t, err)

	var js2 JsonScan[map[string]interface{}]
	err = json.Unmarshal(b, &js2)
	assert.NoError(t, err)
	assert.Equal(t, data["key"], js2.Data()["key"])
}

func TestJsonScan_RoundTrip_Value_Scan(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	js := NewJsonScan(data)

	val, err := js.Value()
	assert.NoError(t, err)

	var js2 JsonScan[map[string]interface{}]
	err = js2.Scan(val)
	assert.NoError(t, err)
	assert.Equal(t, "bar", js2.Data()["foo"])
}

// ======================== Data / Set tests ========================

func TestJsonScan_Data_NonNil(t *testing.T) {
	data := map[string]interface{}{"key": "val"}
	js := NewJsonScan(data)
	result := js.Data()
	assert.Equal(t, data, result)
}

func TestJsonScan_Data_NilV(t *testing.T) {
	var js JsonScan[map[string]interface{}]
	// When v is nil, Data() allocates new(T) and returns the zero value map (empty, non-nil)
	result := js.Data()
	assert.IsType(t, map[string]interface{}(nil), result)
}

func TestJsonScan_Set_UpdatesValue(t *testing.T) {
	js := NewJsonScan(map[string]interface{}{"old": "value"})
	js.Set(map[string]interface{}{"new": "data"})
	// Set replaces v entirely, old key is gone
	assert.Equal(t, map[string]interface{}{"new": "data"}, js.Data())
}
