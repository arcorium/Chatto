package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType uint8

const (
	NotifTyping NotificationType = iota
	NotifOnline
	NotifOffline
	NotifJoinRoom
	NotifLeaveRoom
)

func NewNotificationPayload(notification *Notification) *Payload {
	return &Payload{
		Type:     PayloadNotification,
		ClientId: "",
		Data: NotificationRespond{
			Type:      notification.Type,
			Sender:    "",
			Message:   notification.Message,
			Timestamp: notification.Timestamp,
		},
	}
}

func NewNotification(senderUserId string, receiverId string, notifType NotificationType, message ...string) *Notification {
	notif := &Notification{
		Type:     notifType,
		Receiver: receiverId,
	}
	if message != nil && len(message) > 0 {
		notif.Message = message[0]
	}
	notif.Populate()
	return notif
}

func NewServerNotification(receiverId string, notifType NotificationType, message ...string) *Notification {
	return NewNotification("server", receiverId, notifType, message...)
}

type Notification struct {
	Id        string           `json:"id,omitempty"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message,omitempty"`
	Receiver  string           `json:"receiver"`
	Timestamp int64            `json:"timestamp,omitempty"`
}

// NotificationRespond Used as respond to websocket client
type NotificationRespond struct {
	Type      NotificationType `json:"type"`
	Sender    string           `json:"sender,omitempty"`
	Message   string           `json:"message,omitempty"`
	Timestamp int64            `json:"timestamp"`
}

func (n *Notification) Populate() {
	_, err := uuid.Parse(n.Id)
	if err != nil {
		n.Id = uuid.NewString()
	}

	n.Timestamp = time.Now().Unix()
}
