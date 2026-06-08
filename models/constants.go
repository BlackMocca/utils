package models

import (
	"4d63.com/tz"
)

const (
	TimestampLayout = "2006-01-02 15:04:05"
	DateLayout      = "2006-01-02"
)

var (
	Location, _ = tz.LoadLocation("Asia/Bangkok")
)
