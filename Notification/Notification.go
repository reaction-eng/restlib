package Notification

import (
	"time"
)

type Notification struct {
	Priority   int
	Send       time.Time
	Expiration time.Time
	Message    string
}
