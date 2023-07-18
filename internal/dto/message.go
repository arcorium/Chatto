package dto

import (
	"chatto/internal/model"
	"github.com/google/uuid"
	"time"
)

func NewMessageFromInput(sender *model.Client, message *MessageInput) model.Message {
	return model.Message{
		// Check id
		Id:        uuid.NewString(),
		SenderId:  sender.UserId,
		Message:   message.Message,
		Timestamp: time.Now().Unix(),
	}
}

type MessageInput struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func NewMessageOutput(input *MessageInput, message *model.Message) MessageOutput {
	return MessageOutput{
		Id:         message.Id,
		SenderId:   message.SenderId,
		ReceiverId: input.Receiver,
		Message:    message.Message,
		Timestamp:  message.Timestamp,
	}
}

type MessageOutput struct {
	Id         string `json:"id"`
	SenderId   string `json:"sender"`
	ReceiverId string `json:"receiver"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"ts"`
}

func NewMessageForward(receiverId string, message *MessageOutput) MessageForward {
	return MessageForward{
		ReceiverId: receiverId,
		SenderId:   message.SenderId,
		Message:    message.Message,
		Timestamp:  message.Timestamp,
	}
}

type MessageForward struct {
	ReceiverId string `json:"recv_id"`
	SenderId   string `json:"sender_id"` // Used for testing
	Message    string `json:"message"`
	Timestamp  int64  `json:"ts"`
}
