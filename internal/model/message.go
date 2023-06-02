package model

import (
	"time"

	"github.com/google/uuid"
)

func NewMessagePayload(client *Client, message *Message) *Payload {
	return &Payload{
		Type: PayloadMessage,
		Data: MessageRespond{
			RoomId:    message.Receiver,
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		},
	}
}

type Message struct {
	Id        string `json:"id,omitempty"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type MessageRespond struct {
	RoomId    string `json:"room_id"`
	SenderId  string `json:"sender_id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (m *Message) Populate() {
	// Check id
	_, err := uuid.Parse(m.Id)
	if err != nil {
		m.Id = uuid.NewString()
	}

	m.Timestamp = time.Now().Unix()
}
