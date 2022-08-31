package database

import "time"

type Message struct {
	ID        int
	City      string
	Timestamp time.Time
	Text      string
	Media     string
	Place     string
	From      string
}
