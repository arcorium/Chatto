package model

import (
	"encoding/json"
)

const (
	PayloadForwarder           = "forwarder"
	PayloadPrivateNotification = "private-notif"
	PayloadRoomNotification    = "room-notif"
	PayloadPrivateChat         = "private-chat"
	PayloadRoomChat            = "room-chat"
	PayloadCreateRoom          = "create-room"
	PayloadJoinRoom            = "join-room"
	PayloadLeaveRoom           = "leave-room"
	PayloadGetUsers            = "get-users"
	PayloadError               = "error"
)

func Decode[T any](bytes []byte) (T, error) {
	var t T
	err := json.Unmarshal(bytes, &t)
	return t, err
}

func NewErrorPayload(message string) Payload {
	return Payload{
		Type: PayloadError,
		Data: struct {
			string `json:"message"`
		}{message},
	}
}

func NewPayload[T any](types string, data *T) Payload {
	return Payload{
		Type: types,
		Data: *data,
	}
}

type Payload struct {
	Type     string `json:"type"`
	ClientId string `json:"-"`
	Data     any    `json:"data"`
}

func (p *Payload) EncodeData() ([]byte, error) {
	return json.Marshal(p.Data)
}

func (p *Payload) Populate(client *Client) {
	p.ClientId = client.Id
}

type ForwarderType uint8

const (
	ForwardNotification ForwarderType = iota
	ForwardMessage
)

func NewOutcomeNotificationForward(receiver *Client, notification *OutcomePrivateNotification) OutcomeForward {
	return OutcomeForward{
		ReceiverId:       receiver.UserId,
		Receiver:         receiver.Username,
		SenderId:         notification.Sender,
		Type:             ForwardNotification,
		NotificationType: notification.Type,
		Message:          notification.Message,
		Timestamp:        notification.Timestamp,
	}
}

func NewOutcomeMessageForward(receiver *Client, message *OutcomePrivateMessage) OutcomeForward {
	return OutcomeForward{
		ReceiverId:       receiver.UserId,
		Receiver:         receiver.Username,
		SenderId:         message.Sender,
		Type:             ForwardMessage,
		NotificationType: -1,
		Message:          message.Message,
		Timestamp:        message.Timestamp,
	}
}

type OutcomeForward struct {
	ReceiverId       string           `json:"receiver_id"`
	Receiver         string           `json:"receiver"`
	SenderId         string           `json:"sender_id"` // Used for testing
	Type             ForwarderType    `json:"type"`
	NotificationType NotificationType `json:"notification_type"` // Only used when the forwarded payload is notification
	Message          string           `json:"message"`
	Timestamp        int64            `json:"ts"`
}
