package model

type MessageType string

type Message struct {
	Id         string `json:"id"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"ts"`
}
