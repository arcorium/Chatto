package model

import (
	"time"

	"github.com/google/uuid"
)

func NewMessagePayload(client *Client, message *Message) *Payload {
	return &Payload{
		Type: PayloadMessage,
		Data: MessageRespond{
			Sender:    client.UserId,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		},
	}
}

func NewMessage(receiver string, message string) *Message {
	msg := &Message{
		Receiver: receiver,
		Message:  message,
	}
	msg.Populate()
	return msg
}

type Message struct {
	Id        string `json:"id,omitempty"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type MessageRespond struct {
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
