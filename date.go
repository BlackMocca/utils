package utils

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	DateLayout = "2006-01-02"
)

type Date time.Time

/*
------------------------
Date Function
------------------------
*/

func NewDateFromNow() Date {
	return NewDateFromString(time.Now().Format(DateLayout))
}

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

func (j *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
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

func (j Date) String() string {
	return j.Format(DateLayout)
}

func (j Date) ToPointer() *Date {
	return &j
}
