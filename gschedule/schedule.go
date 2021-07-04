package gscheduler

import (
	"time"
)

type Schedule interface {
	Next(t time.Time) time.Time
}
