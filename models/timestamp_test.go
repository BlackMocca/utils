package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ======================== MarshalJSON tests ========================

func TestTimestamp_MarshalJSON_Valid(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	b, err := json.Marshal(ts)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-06-15 14:30:45"`, string(b))
}

func TestTimestamp_MarshalJSON_Zero(t *testing.T) {
	ts := Timestamp{}
	b, err := json.Marshal(ts)
	assert.NoError(t, err)
	assert.Equal(t, `"0001-01-01 00:00:00"`, string(b))
}

func TestTimestamp_MarshalJSON_WithMicroseconds(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 123456000, Location))
	b, err := json.Marshal(ts)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-06-15 14:30:45"`, string(b))
}

func TestTimestamp_MarshalJSON_DifferentTimes(t *testing.T) {
	tests := []struct {
		name     string
		hour     int
		minute   int
		second   int
		expected string
	}{
		{"Midnight", 0, 0, 0, `"2025-06-15 00:00:00"`},
		{"Noon", 12, 0, 0, `"2025-06-15 12:00:00"`},
		{"End of day", 23, 59, 59, `"2025-06-15 23:59:59"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := Timestamp(time.Date(2025, 6, 15, tc.hour, tc.minute, tc.second, 0, Location))
			b, err := json.Marshal(ts)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(b))
		})
	}
}

// ======================== UnmarshalJSON tests ========================

func TestTimestamp_UnmarshalJSON_Standard(t *testing.T) {
	var ts Timestamp
	input := []byte(`"2025-06-15 14:30:45"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}

func TestTimestamp_UnmarshalJSON_Null(t *testing.T) {
	var ts Timestamp
	input := []byte(`null`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
}

func TestTimestamp_UnmarshalJSON_EmptyString(t *testing.T) {
	var ts Timestamp
	input := []byte(`""`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
}

func TestTimestamp_UnmarshalJSON_NilString(t *testing.T) {
	var ts Timestamp
	input := []byte(`"nil"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
}

func TestTimestamp_UnmarshalJSON_Microseconds(t *testing.T) {
	// Input: 12:05:12 UTC → code does t.Add(-7h) = 05:05:12 UTC → In(Location) = 12:05:12 ICT
	var ts Timestamp
	input := []byte(`"2025-04-03 12:05:12.510131"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "2025-04-03 12:05:12", ts.String())
}

func TestTimestamp_UnmarshalJSON_Milliseconds(t *testing.T) {
	var ts Timestamp
	input := []byte(`"2025-04-03 12:05:12.510"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "2025-04-03 12:05:12", ts.String())
}

func TestTimestamp_UnmarshalJSON_DateOnly(t *testing.T) {
	// Date-only format: parsed with no time component, stored as-is
	var ts Timestamp
	input := []byte(`"2024-06-25"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// The layout "2006-01-02" parses with fixed zone location (not UTC),
	// so the UTC offset conversion isn't applied
	assert.Equal(t, "2024-06-25 00:00:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_ISO8601(t *testing.T) {
	// ISO8601 without timezone: parsed with fixed zone location
	var ts Timestamp
	input := []byte(`"2024-06-25T14:30:00"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// The layout "2006-01-02T15:04:05" parses with fixed zone (not UTC),
	// so no offset conversion is applied
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_RFC3339(t *testing.T) {
	// RFC3339 includes Z timezone → parsed as UTC with location UTC
	// Code does t.Add(-7h) then In(Location)
	var ts Timestamp
	input := []byte(`"2025-04-03T12:05:12Z"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// 12:05:12 UTC → subtract 7h = 05:05:12 UTC → convert to Bangkok = 12:05:12 ICT
	assert.Equal(t, "2025-04-03 12:05:12", ts.String())
}

func TestTimestamp_UnmarshalJSON_RFC3339Nano(t *testing.T) {
	var ts Timestamp
	input := []byte(`"2025-04-03T12:05:12.510131Z"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "2025-04-03 12:05:12", ts.String())
}

func TestTimestamp_UnmarshalJSON_RFC1123(t *testing.T) {
	var ts Timestamp
	input := []byte(`"Mon, 25 Jun 2024 14:30:00 UTC"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// 14:30:00 UTC → -7h = 07:30:00 UTC → In(Location) = 14:30:00 ICT
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_RFC1123Z(t *testing.T) {
	var ts Timestamp
	input := []byte(`"Mon, 25 Jun 2024 14:30:00 +0000"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// +0000 is UTC → same conversion as above
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_RFC850(t *testing.T) {
	var ts Timestamp
	input := []byte(`"Monday, 25-Jun-24 14:30:00 UTC"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_ANSIC(t *testing.T) {
	// No timezone info → location is UTC (Go default for time.Parse with no TZ in format)
	var ts Timestamp
	input := []byte(`"Mon Jun 25 14:30:00 2024"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// 14:30:00 UTC → -7h = 07:30:00 UTC → In(Location) = 14:30:00 ICT
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_UnixDate(t *testing.T) {
	var ts Timestamp
	input := []byte(`"Mon Jun 25 14:30:00 MST 2024"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// MST parsed without timezone → location is UTC
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_RubyDate(t *testing.T) {
	var ts Timestamp
	input := []byte(`"Mon Jun 25 14:30:00 -0700 2024"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// -0700 → location is fixed zone UTC-7, not UTC, so no conversion applied
	assert.Equal(t, "2024-06-25 14:30:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_DateOnlyNoTime(t *testing.T) {
	var ts Timestamp
	input := []byte(`"25 Jun 2024"`)
	err := json.Unmarshal(input, &ts)
	assert.NoError(t, err)
	// No timezone info → no UTC conversion applied
	assert.Equal(t, "2024-06-25 00:00:00", ts.String())
}

func TestTimestamp_UnmarshalJSON_InvalidFormat(t *testing.T) {
	tests := []string{`"not-a-timestamp"`, `"2025-13-01"`, `"invalid"`}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var ts Timestamp
			err := json.Unmarshal([]byte(input), &ts)
			assert.Error(t, err)
		})
	}
}

// ======================== Scan tests ========================

func TestTimestamp_Scan_TimeTime_UTC(t *testing.T) {
	var ts Timestamp
	input := time.Date(2025, 6, 15, 14, 30, 45, 0, time.UTC)
	err := ts.Scan(input)
	assert.NoError(t, err)

	// Scan also applies the same UTC → Bangkok conversion logic
	// 14:30:45 UTC → -7h = 07:30:45 UTC → In(Location) = 14:30:45 ICT
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}

func TestTimestamp_Scan_TimeTime_Local(t *testing.T) {
	var ts Timestamp
	input := time.Date(2025, 6, 15, 14, 30, 45, 0, Location)
	err := ts.Scan(input)
	assert.NoError(t, err)

	// Not UTC → no conversion applied
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}

func TestTimestamp_Scan_String(t *testing.T) {
	var ts Timestamp
	input := "2025-06-15 14:30:45"
	err := ts.Scan(input)
	assert.NoError(t, err)
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}

func TestTimestamp_Scan_Bytes(t *testing.T) {
	var ts Timestamp
	input := []byte("2025-06-15 14:30:45")
	err := ts.Scan(input)
	assert.NoError(t, err)
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}

func TestTimestamp_Scan_Nil(t *testing.T) {
	var ts Timestamp
	err := ts.Scan(nil)
	assert.NoError(t, err)
	assert.True(t, ts.ToTime().IsZero())
}

func TestTimestamp_Scan_InvalidType(t *testing.T) {
	var ts Timestamp
	input := int(42)
	err := ts.Scan(input)
	assert.Error(t, err)
	assert.Equal(t, "cannot scan type int into Timestamp", err.Error())
}

func TestTimestamp_Scan_InvalidTimeString(t *testing.T) {
	var ts Timestamp
	input := "not-a-timestamp"
	err := ts.Scan(input)
	assert.Error(t, err)
}

func TestTimestamp_Scan_DifferentDates(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-01-01 00:00:00", "2024-01-01 00:00:00"},
		{"2024-12-31 23:59:59", "2024-12-31 23:59:59"},
		{"2024-06-15 12:30:45", "2024-06-15 12:30:45"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			var ts Timestamp
			err := ts.Scan([]byte(tc.input))
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, ts.String())
		})
	}
}

// ======================== Value tests ========================

func TestTimestamp_Value_Valid(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	val, err := ts.Value()
	assert.NoError(t, err)
	assert.Equal(t, "2025-06-15 14:30:45", val)
}

func TestTimestamp_Value_Zero(t *testing.T) {
	var ts Timestamp
	val, err := ts.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestTimestamp_Value_DifferentTimes(t *testing.T) {
	tests := []string{
		"2024-01-01 00:00:00",
		"2024-12-31 23:59:59",
		"2024-06-15 12:30:45",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			ts := NewTimestampFromString(tc)
			val, err := ts.Value()
			assert.NoError(t, err)
			assert.Equal(t, tc, val)
		})
	}
}

func TestTimestamp_ValueOrZero_Valid(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	result := ts.ValueOrZero()
	assert.Equal(t, "2025-06-15 14:30:45", result)
}

func TestTimestamp_ValueOrZero_Zero(t *testing.T) {
	var ts Timestamp
	result := ts.ValueOrZero()
	assert.Equal(t, "", result)
}

// ======================== Format / YearDay tests ========================

func TestTimestamp_Format(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))

	tests := []struct {
		format   string
		expected string
	}{
		{"2006-01-02", "2025-06-15"},
		{"Mon Jan 2 2006 15:04:05 MST", "Sun Jun 15 2025 14:30:45 +07"},
		{"2006/01/02 15:04:05", "2025/06/15 14:30:45"},
	}

	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {
			result := ts.Format(tc.format)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTimestamp_YearDay(t *testing.T) {
	ts := Timestamp(time.Date(2025, 1, 1, 0, 0, 0, 0, Location))
	assert.Equal(t, 1, ts.YearDay())

	ts = Timestamp(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))
	assert.Equal(t, 166, ts.YearDay())

	ts = Timestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, Location))
	assert.Equal(t, 366, ts.YearDay())
}

// ======================== ToUnix tests ========================

func TestTimestamp_ToUnix(t *testing.T) {
	ts := Timestamp(time.Unix(1750080000, 0).In(Location))
	unix := ts.ToUnix()
	assert.Equal(t, int64(1750080000), unix)
}

// ======================== Round-trip tests ========================

func TestTimestamp_RoundTrip_MarshalJSON_UnmarshalJSON(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	var ts2 Timestamp
	err = json.Unmarshal(b, &ts2)
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}

func TestTimestamp_RoundTrip_Value_Scan_String(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	val, err := ts.Value()
	assert.NoError(t, err)

	var ts2 Timestamp
	err = ts2.Scan(val.(string))
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}

func TestTimestamp_RoundTrip_Value_Scan(t *testing.T) {
	ts := Timestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, Location))
	val, err := ts.Value()
	assert.NoError(t, err)

	var ts2 Timestamp
	if val != nil {
		err = ts2.Scan(val.(string))
	} else {
		err = ts2.Scan(nil)
	}
	assert.NoError(t, err)
	assert.Equal(t, ts.String(), ts2.String())
}

func TestTimestamp_String(t *testing.T) {
	ts := Timestamp(time.Date(2025, 6, 15, 14, 30, 45, 0, Location))
	assert.Equal(t, "2025-06-15 14:30:45", ts.String())
}
