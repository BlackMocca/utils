package models

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ======================== MarshalJSON tests ========================

func TestDate_MarshalJSON_Valid(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))
	b, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-06-15"`, string(b))
}

func TestDate_MarshalJSON_Zero(t *testing.T) {
	d := Date{}
	b, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"0001-01-01"`, string(b))
}

func TestDate_MarshalJSON_WithTimeComponents(t *testing.T) {
	// Date should ignore time components and use only year/month/day in Location timezone
	d := Date(time.Date(2025, 6, 15, 14, 30, 45, 999000000, Location))
	b, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-06-15"`, string(b))
}

func TestDate_MarshalJSON_DifferentDates(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"Jan 1", time.Date(2024, time.January, 1, 12, 0, 0, 0, Location), `"2024-01-01"`},
		{"Dec 31", time.Date(2099, time.December, 31, 12, 0, 0, 0, Location), `"2099-12-31"`},
		{"Feb 29 leap", time.Date(2024, time.February, 29, 12, 0, 0, 0, Location), `"2024-02-29"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := Date(tc.input)
			b, err := json.Marshal(d)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(b))
		})
	}
}

// ======================== UnmarshalJSON tests ========================

func TestDate_UnmarshalJSON_Valid(t *testing.T) {
	var d Date
	input := []byte(`"2025-06-15"`)
	err := json.Unmarshal(input, &d)
	assert.NoError(t, err)

	yr, mo, dy := time.Time(d).Date()
	assert.Equal(t, 2025, yr)
	assert.Equal(t, time.Month(6), mo)
	assert.Equal(t, 15, dy)
}

func TestDate_UnmarshalJSON_Null(t *testing.T) {
	var d Date
	input := []byte(`null`)
	err := json.Unmarshal(input, &d)
	assert.NoError(t, err)
}

func TestDate_UnmarshalJSON_EmptyString(t *testing.T) {
	var d Date
	input := []byte(`""`)
	err := json.Unmarshal(input, &d)
	assert.NoError(t, err)
}

func TestDate_UnmarshalJSON_NilString(t *testing.T) {
	var d Date
	input := []byte(`"nil"`)
	err := json.Unmarshal(input, &d)
	assert.NoError(t, err)
}

func TestDate_UnmarshalJSON_InvalidFormat(t *testing.T) {
	tests := []string{`"invalid"`, `"2025/06/15"`, `"15-06-2025"`, `"2025-13-01"`, `"2025-06-32"`}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var d Date
			err := json.Unmarshal([]byte(input), &d)
			assert.Error(t, err)
		})
	}
}

func TestDate_UnmarshalJSON_ValidDates(t *testing.T) {
	tests := []struct {
		input    string
		year     int
		month    time.Month
		day      int
	}{
		{`"2024-01-01"`, 2024, time.January, 1},
		{`"2024-12-31"`, 2024, time.December, 31},
		{`"2024-02-29"`, 2024, time.February, 29}, // leap year
		{`"0001-01-01"`, 1, time.January, 1},       // zero date
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			var d Date
			err := json.Unmarshal([]byte(tc.input), &d)
			assert.NoError(t, err)
			yr, mo, dy := time.Time(d).Date()
			assert.Equal(t, tc.year, yr)
			assert.Equal(t, tc.month, mo)
			assert.Equal(t, tc.day, dy)
		})
	}
}

// ======================== Scan tests ========================

func TestDate_Scan_TimeTime(t *testing.T) {
	var d Date
	input := time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC)
	err := d.Scan(input)
	assert.NoError(t, err)

	yr, mo, dy := time.Time(d).Date()
	assert.Equal(t, 2025, yr)
	assert.Equal(t, time.Month(6), mo)
	assert.Equal(t, 15, dy)
}

func TestDate_Scan_String(t *testing.T) {
	var d Date
	input := "2025-06-15"
	err := d.Scan(input)
	assert.NoError(t, err)

	yr, mo, dy := time.Time(d).Date()
	assert.Equal(t, 2025, yr)
	assert.Equal(t, time.Month(6), mo)
	assert.Equal(t, 15, dy)
}

func TestDate_Scan_Bytes(t *testing.T) {
	var d Date
	input := []byte("2025-06-15")
	err := d.Scan(input)
	assert.NoError(t, err)

	yr, mo, dy := time.Time(d).Date()
	assert.Equal(t, 2025, yr)
	assert.Equal(t, time.Month(6), mo)
	assert.Equal(t, 15, dy)
}

func TestDate_Scan_Nil(t *testing.T) {
	var d Date
	err := d.Scan(nil)
	assert.NoError(t, err)
	assert.True(t, time.Time(d).IsZero())
}

func TestDate_Scan_InvalidType(t *testing.T) {
	var d Date
	input := int(42)
	err := d.Scan(input)
	assert.Error(t, err)
	assert.Equal(t, "cannot scan type int into Date", err.Error())
}

func TestDate_Scan_InvalidDateString(t *testing.T) {
	var d Date
	input := "not-a-date"
	err := d.Scan(input)
	assert.Error(t, err)
}

func TestDate_Scan_DifferentDates(t *testing.T) {
	tests := []string{"2024-01-01", "2024-12-31", "2024-02-29", "0001-01-01"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var d Date
			err := d.Scan([]byte(input))
			assert.NoError(t, err)
			assert.Equal(t, input, d.String())
		})
	}
}

// ======================== Value tests ========================

func TestDate_Value_Valid(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))
	val, err := d.Value()
	assert.NoError(t, err)
	assert.Equal(t, "2025-06-15", val)
}

func TestDate_Value_Zero(t *testing.T) {
	var d Date
	val, err := d.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestDate_Value_DifferentDates(t *testing.T) {
	tests := []struct {
		input    string
		expected driver.Value
	}{
		{"2024-01-01", "2024-01-01"},
		{"2024-12-31", "2024-12-31"},
		{"2024-02-29", "2024-02-29"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			d := NewDateFromString(tc.input)
			val, err := d.Value()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, val)
		})
	}
}

// ======================== Round-trip tests ========================

func TestDate_RoundTrip_MarshalJSON_UnmarshalJSON(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 14, 30, 45, 999000000, Location))
	b, err := json.Marshal(d)
	assert.NoError(t, err)

	var d2 Date
	err = json.Unmarshal(b, &d2)
	assert.NoError(t, err)
	assert.Equal(t, d.String(), d2.String())
}

func TestDate_RoundTrip_Value_Scan_String(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))
	val, err := d.Value()
	assert.NoError(t, err)

	var d2 Date
	err = d2.Scan(val.(string))
	assert.NoError(t, err)
	assert.Equal(t, d.String(), d2.String())
}

func TestDate_RoundTrip_Value_Scan(t *testing.T) {
	d := Date(time.Date(2024, 12, 31, 0, 0, 0, 0, Location))
	val, err := d.Value()
	assert.NoError(t, err)

	var d2 Date
	if val != nil {
		err = d2.Scan(val.(string)) // Value() returns string for Date
	} else {
		err = d2.Scan(nil)
	}
	assert.NoError(t, err)
	assert.Equal(t, d.String(), d2.String())
}

func TestDate_RoundTrip_Value_Scan_Time(t *testing.T) {
	d := Date(time.Date(2025, 1, 1, 0, 0, 0, 0, Location))
	val, err := d.Value()
	assert.NoError(t, err)

	var d2 Date
	err = d2.Scan(val.(string)) // Value returns string, not time.Time
	assert.NoError(t, err)
	assert.Equal(t, d.String(), d2.String())
}

// ======================== String / Format tests ========================

func TestDate_String(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))
	assert.Equal(t, "2025-06-15", d.String())
}

func TestDate_Format(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location))

	tests := []struct {
		format   string
		expected string
	}{
		{"2006-01-02", "2025-06-15"},
		{"Mon Jan 2 2006", "Sun Jun 15 2025"},
		{"2006/01/02", "2025/06/15"},
	}

	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {
			result := d.Format(tc.format)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDate_Weekday(t *testing.T) {
	d := Date(time.Date(2025, 6, 15, 0, 0, 0, 0, Location)) // Sunday
	assert.Equal(t, time.Sunday, d.Weekday())

	d2 := Date(time.Date(2025, 6, 16, 0, 0, 0, 0, Location)) // Monday
	assert.Equal(t, time.Monday, d2.Weekday())
}
