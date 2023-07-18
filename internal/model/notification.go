package model

type NotificationType int8

const (
	NotifTyping NotificationType = iota
	NotifOnline
	NotifOffline
	NotifJoinRoom
	NotifLeaveRoom
)

type Notification struct {
	Id        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	SenderId  string           `json:"sender"`
	Timestamp int64            `json:"ts"`
}
