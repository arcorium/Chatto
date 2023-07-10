package model

import (
	"time"

	"github.com/google/uuid"
)

func NewRoomMessagePayload(client *Client, message *Message) Payload {
	return Payload{
		Type: PayloadRoomChat,
		Data: OutcomeRoomMessage{
			RoomId:    message.Receiver,
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		},
	}
}

func NewMessage(sender *Client, message *IncomeMessage) *Message {
	return &Message{
		// Check id
		Id:        uuid.NewString(),
		Sender:    sender.UserId,
		Receiver:  message.Receiver,
		Message:   message.Message,
		Timestamp: time.Now().Unix(),
	}
}

type Message struct {
	Id        string `json:"id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}

type IncomeMessage struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func NewOutcomePrivateMessage(sender *Client, message *Message) OutcomePrivateMessage {
	return OutcomePrivateMessage{
		SenderId:  sender.UserId,
		Sender:    sender.Username,
		Message:   message.Message,
		Timestamp: message.Timestamp,
	}
}

type OutcomePrivateMessage struct {
	SenderId  string `json:"sender_id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}

func NewOutcomeRoomMessage(sender *Client, message *Message) OutcomeRoomMessage {
	return OutcomeRoomMessage{
		RoomId:    message.Receiver,
		SenderId:  sender.UserId,
		Sender:    sender.Username,
		Message:   message.Message,
		Timestamp: message.Timestamp,
	}
}

type OutcomeRoomMessage struct {
	RoomId    string `json:"room_id"`
	SenderId  string `json:"sender_id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}
