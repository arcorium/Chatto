package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType int8

const (
	NotifTyping NotificationType = iota
	NotifOnline
	NotifOffline
	NotifInvitedRoom
	NotifJoinRoom
	NotifLeaveRoom
)

func NewNotificationPayload(client *Client, notification *Notification) *Payload {
	return &Payload{
		Type:     PayloadPrivateNotification,
		ClientId: client.Id,
		Data: OutcomePrivateNotification{
			Type:      notification.Type,
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   notification.Message,
			Timestamp: notification.Timestamp,
		},
	}
}

func NewNotification(notification *IncomeNotification) *Notification {
	return &Notification{
		Id:        uuid.NewString(),
		Type:      notification.Type,
		Message:   "", // TODO: Set message based on notification type
		Receiver:  notification.Receiver,
		Timestamp: time.Now().Unix(),
	}
}

type IncomeNotification struct {
	Type     NotificationType `json:"type"`
	Receiver string           `json:"receiver"`
}

type Notification struct {
	Id        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	Sender    string           `json:"sender"`
	Receiver  string           `json:"receiver"`
	Timestamp int64            `json:"ts"`
}

func NewOutcomeNotification(sender *Client, notification *Notification) OutcomePrivateNotification {
	return OutcomePrivateNotification{
		Type:      notification.Type,
		SenderId:  sender.UserId,
		Sender:    sender.Username,
		Message:   notification.Message,
		Timestamp: notification.Timestamp,
	}
}

type OutcomePrivateNotification struct {
	Type      NotificationType `json:"type"`
	SenderId  string           `json:"sender_id"`
	Sender    string           `json:"sender"`
	Message   string           `json:"message"`
	Timestamp int64            `json:"ts"`
}

func NewOutcomeRoomNotification(sender *Client, notification *Notification) OutcomeRoomNotification {
	return OutcomeRoomNotification{
		Type:      notification.Type,
		RoomId:    notification.Receiver,
		SenderId:  sender.UserId,
		Sender:    sender.Username,
		Message:   notification.Message,
		Timestamp: notification.Timestamp,
	}
}

type OutcomeRoomNotification struct {
	Type      NotificationType `json:"type"`
	RoomId    string           `json:"room_id"`
	SenderId  string           `json:"sender_id"`
	Sender    string           `json:"sender"`
	Message   string           `json:"message"`
	Timestamp int64            `json:"ts"`
}
