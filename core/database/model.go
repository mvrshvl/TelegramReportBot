package database

import "time"

type Message struct {
	Timestamp time.Time
	Text      string
	Media     string
	Place     string
	From      string
}
