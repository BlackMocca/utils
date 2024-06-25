package utils

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"4d63.com/tz"
)

const (
	TimestampLayout = "2006-01-02 15:04:05"
)

var (
	Location, _ = tz.LoadLocation("Asia/Bangkok")
)

type Timestamp time.Time

/*
------------------------
Timestamp Function
------------------------
*/

func NewTimestampFromNow() Timestamp {
	return NewTimestampFromTime(time.Now())
}

func NewTimestampFromString(dateString string) Timestamp {
	if dateString == "" {
		return Timestamp(time.Time{})
	}
	d, err := time.ParseInLocation(TimestampLayout, dateString, Location)
	if err != nil {
		panic(err)
	}
	return Timestamp(d.In(Location))
}

func NewTimestampFromTime(t time.Time) Timestamp {
	d, err := time.ParseInLocation(TimestampLayout, t.In(Location).Format(TimestampLayout), Location)
	if err != nil {
		panic(err)
	}

	return Timestamp(d.In(Location))
}

func (t Timestamp) ToUnix() int64 {
	tt := time.Time(t)
	return tt.Unix()
}

func (j *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(TimestampLayout, s)
	if err != nil {
		return err
	}
	*j = Timestamp(t)
	return nil
}

func (j Timestamp) ToPointer() *Timestamp {
	return &j
}

func (j Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Format(TimestampLayout))
}

func (j Timestamp) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func (j Timestamp) String() string {
	return j.Format(TimestampLayout)
}

func (j Timestamp) Value() (driver.Value, error) {
	if j == (Timestamp{}) {
		return nil, nil
	}
	return j.String(), nil
}

func (j Timestamp) ToTime() time.Time {
	return time.Time(j)
}
