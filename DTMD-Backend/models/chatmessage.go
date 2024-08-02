package models

import "time"

type ChatMessage struct {
	Sender    string    `json:"sender"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
