package dto

import (
	"chatto/internal/model"
	"github.com/google/uuid"
	"time"
)

func NewNotificationOutput(senderId string, input *NotificationInput, notification *model.Notification) NotificationOutput {
	return NotificationOutput{
		Type:       notification.Type,
		SenderId:   senderId,
		ReceiverId: input.Receiver,
		Message:    notification.Message,
		Timestamp:  notification.Timestamp,
	}
}

type NotificationOutput struct {
	Type       model.NotificationType `json:"type"`
	SenderId   string                 `json:"sender_id"`
	ReceiverId string                 `json:"recv_id"`
	Message    string                 `json:"message"`
	Timestamp  int64                  `json:"ts"`
}

func NewNotificationFromInput(senderId string, notification *NotificationInput) model.Notification {
	return model.Notification{
		Id:        uuid.NewString(),
		Type:      notification.Type,
		Message:   "", // TODO: Set message based on notification type
		Timestamp: time.Now().Unix(),
		SenderId:  senderId,
	}
}

type NotificationInput struct {
	Type     model.NotificationType `json:"type"`
	Receiver string                 `json:"receiver"`
}

func NewNotificationForward(receiverId string, notification *NotificationOutput) NotificationForward {
	return NotificationForward{
		ReceiverId:       receiverId,
		SenderId:         notification.SenderId,
		NotificationType: notification.Type,
		Message:          notification.Message,
		Timestamp:        notification.Timestamp,
	}
}

type NotificationForward struct {
	ReceiverId       string                 `json:"receiver_id"`
	SenderId         string                 `json:"sender_id"`         // Used for testing
	NotificationType model.NotificationType `json:"notification_type"` // Only used when the forwarded payload is notification
	Message          string                 `json:"message"`
	Timestamp        int64                  `json:"ts"`
}
