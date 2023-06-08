package model

import (
	"github.com/google/uuid"
	"time"
)

func NewPrivateMessagePayload(client *Client, message *Message) Payload {
	return Payload{
		Type: PayloadPrivateChat,
		Data: PrivateMessage{
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		},
	}
}

func NewRoomMessagePayload(client *Client, message *Message) Payload {
	return Payload{
		Type: PayloadRoomChat,
		Data: RoomMessage{
			RoomId:    message.Receiver,
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		},
	}
}

func NewMessage(message *IncomeMessage) *Message {
	return &Message{
		// Check id
		Id:        uuid.NewString(),
		Receiver:  message.Receiver,
		Message:   message.Message,
		Timestamp: time.Now().Unix(),
	}
}

type Message struct {
	Id        string `json:"id"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}

type IncomeMessage struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type PrivateMessage struct {
	SenderId  string `json:"sender_id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}

type RoomMessage struct {
	RoomId    string `json:"room_id"`
	SenderId  string `json:"sender_id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}
