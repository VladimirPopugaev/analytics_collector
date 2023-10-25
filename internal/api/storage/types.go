package storage

import "time"

type UserActionInfo struct {
	Time   time.Time
	UserID string
	Data   Request
}

type Request struct {
	Body    string
	Headers string
}

//TODO: add own errors
