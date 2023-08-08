package model

type NotificationType int8

const (
	NotifTyping NotificationType = iota
	NotifJoinRoom
	NotifLeaveRoom
)

func GetNotificationMessage(client *Client, types NotificationType) string {
	message := ""
	switch types {
	case NotifJoinRoom:
		message = client.Username + " joined room"
	case NotifLeaveRoom:
		message = client.Username + " leave room"
	}
	return message
}

type Notification struct {
	Id   string           `json:"id"`
	Type NotificationType `json:"type"`
	//Message   string           `json:"message"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Timestamp  int64  `json:"ts"`
}
