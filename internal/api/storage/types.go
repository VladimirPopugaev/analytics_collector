package storage

import "time"

type UserActionInfo struct {
	Time   time.Time
	UserID string
	Data   RequestInfo
}

type RequestInfo struct {
	Body    string
	Headers string
}
