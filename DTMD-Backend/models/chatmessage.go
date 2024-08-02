package models

type ChatMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
