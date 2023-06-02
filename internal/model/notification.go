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

func NewNotificationPayload(client *Client, notification *Notification) *Payload {
	return &Payload{
		Type:     PayloadNotification,
		ClientId: client.Id,
		Data: NotificationRespond{
			Type:      notification.Type,
			SenderId:  client.UserId,
			Sender:    client.Username,
			Message:   notification.Message,
			Timestamp: notification.Timestamp,
		},
	}
}

func NewNotification(receiverId string, notifType NotificationType, message ...string) *Notification {
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
	SenderId  string           `json:"sender_id,omitempty"`
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
