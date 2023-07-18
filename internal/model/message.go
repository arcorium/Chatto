package model

type Message struct {
	Id        string `json:"id"`
	SenderId  string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}
