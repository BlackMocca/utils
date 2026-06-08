package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date time.Time

/*
------------------------
Date Function
------------------------
*/

func NewDateFromString(dateString string) Date {
	d, err := time.ParseInLocation(DateLayout, dateString, Location)
	if err != nil {
		panic(err)
	}
	return Date(d.In(Location))
}

func NewDateFromTime(t time.Time) Date {
	d, err := time.ParseInLocation(DateLayout, t.Format(DateLayout), Location)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

func ParseDateFromString(dateString string) (Date, error) {
	d, err := time.ParseInLocation(DateLayout, dateString, Location)
	if err != nil {
		return Date(time.Time{}), err
	}
	return Date(d), nil
}

func ParseDateFromTime(t time.Time) (Date, error) {
	d, err := time.ParseInLocation(DateLayout, t.Format(DateLayout), Location)
	if err != nil {
		return Date(time.Time{}), err
	}
	return Date(d), nil
}

func (j *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if len(b) == 0 || string(b) == "null" || s == "" || s == "nil" {
		return nil
	}

	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return err
	}
	*j = Date(t)
	return nil
}

func (j Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Format(DateLayout))
}

func (j Date) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func (j Date) Weekday() time.Weekday {
	t := time.Time(j)
	return t.Weekday()
}

func (j Date) String() string {
	return j.Format(DateLayout)
}

func (j Date) ToTime() time.Time {
	return time.Time(j)
}

func (j Date) ToPointer() *Date {
	return &j
}

func (j Date) ToTimestamp() Timestamp {
	return Timestamp(time.Time(j))
}

func (j *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*j = Date(v)
	case string:
		parsedTime, err := time.ParseInLocation(DateLayout, v, Location)
		if err != nil {
			return err
		}
		*j = Date(parsedTime)
	case []byte:
		parsedTime, err := time.ParseInLocation(DateLayout, string(v), Location)
		if err != nil {
			return err
		}
		*j = Date(parsedTime)
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}

	return nil
}

func (j Date) Value() (driver.Value, error) {
	if j == (Date{}) {
		return nil, nil
	}
	return j.String(), nil
}
