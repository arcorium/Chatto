package dto

import (
	"time"

	"chatto/internal/model"
	"github.com/google/uuid"
)

func NewNotificationOutput(notification *model.Notification) NotificationOutput {
	return NotificationOutput{
		Type:       notification.Type,
		SenderId:   notification.SenderId,
		ReceiverId: notification.ReceiverId,
		//Message:    notification.Message,
		Timestamp: notification.Timestamp,
	}
}

type NotificationOutput struct {
	Type       model.NotificationType `json:"type"`
	SenderId   string                 `json:"sender_id"`
	ReceiverId string                 `json:"receiver_id"`
	//Message    string                 `json:"message"`
	Timestamp int64 `json:"ts"`
}

func NewNotificationFromInput(userId string, input *NotificationInput) model.Notification {
	return model.Notification{
		Id:   uuid.NewString(),
		Type: input.Type,
		//Message:   model.GetNotificationMessage(sender, input.Type),
		Timestamp:  time.Now().Unix(),
		SenderId:   userId,
		ReceiverId: input.ReceiverId,
	}
}

type NotificationInput struct {
	Type model.NotificationType `json:"type"`
	// ReceiverId can be User or Room, handler distinguish it by the payload type
	ReceiverId string `json:"room_id"`
}

type NotificationRequest struct {
	RoomId   string `json:"room_id"`
	FromTime int64  `json:"from_time"`
	ToTime   int64  `json:"to_time"`
}

func NewNotificationResponse(notification *model.Notification) NotificationResponse {
	return NotificationResponse{
		SenderId:  notification.SenderId,
		Type:      notification.Type,
		Timestamp: notification.Timestamp,
	}
}

type NotificationResponse struct {
	SenderId  string                 `json:"sender_id"`
	Type      model.NotificationType `json:"type"`
	Timestamp int64                  `json:"ts"`
}

type TypingInput struct {
	RoomId string `json:"room_id"`
}
