package dto

import (
	"time"

	"chatto/internal/model"
	"github.com/google/uuid"
)

func NewMessageFromInput(sender *model.Client, message *MessageInput) model.Message {
	return model.Message{
		// Check id
		Id:         uuid.NewString(),
		SenderId:   sender.UserId,
		ReceiverId: message.ReceiverId,
		Message:    message.Message,
		Timestamp:  time.Now().Unix(),
	}
}

type MessageInput struct {
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
}

func NewMessageOutput(message *model.Message) MessageOutput {
	return MessageOutput{
		Id:         message.Id,
		SenderId:   message.SenderId,
		ReceiverId: message.ReceiverId,
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

func NewMessageForward(receiverId string, message *MessageOutput) ChatForward {
	return ChatForward{
		ReceiverId: receiverId,
		SenderId:   message.SenderId,
		Message:    message.Message,
		Timestamp:  message.Timestamp,
	}
}

type ChatForward struct {
	ReceiverId string `json:"recv_id"`
	SenderId   string `json:"sender_id"` // Used for testing
	Message    string `json:"message"`
	Timestamp  int64  `json:"ts"`
}

type MessageRequest struct {
	RoomId   string `json:"room_id"`
	FromTime int64  `json:"from_time"`
	ToTime   int64  `json:"to_time"`
}

func NewMessageResponse(message *model.Message) MessageResponse {
	return MessageResponse{
		SenderId:  message.SenderId,
		Message:   message.Message,
		Timestamp: message.Timestamp,
	}
}

type MessageResponse struct {
	SenderId  string `json:"sender_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"ts"`
}
