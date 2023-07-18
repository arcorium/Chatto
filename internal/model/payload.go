package model

import (
	"encoding/json"
)

const (
	PayloadMessageForwarder      = "msg-forward"
	PayloadNotificationForwarder = "notif-forward"
	PayloadPrivateNotification   = "private-notif"
	PayloadRoomNotification      = "room-notif"
	PayloadPrivateChat           = "private-chat"
	PayloadRoomChat              = "room-chat"
	PayloadCreateRoom            = "create-room"
	PayloadJoinRoom              = "join-room"
	PayloadLeaveRoom             = "leave-room"
	PayloadInviteToRoom          = "invite-room"
	PayloadGetUsers              = "get-users"
	PayloadErrorResponse         = "error"
	PayloadSuccessResponse       = "success"
)

func Decode[T any](bytes []byte) (T, error) {
	var t T
	err := json.Unmarshal(bytes, &t)
	return t, err
}

func NewErrorPayloadOutput(code uint, message string) PayloadOutput {
	return PayloadOutput{
		Type: PayloadErrorResponse,
		Data: ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
}

func NewPayloadOutput[T any](types string, data *T) PayloadOutput {
	return PayloadOutput{
		Type: types,
		Data: *data,
	}
}

func NewPayload(client *Client, input *PayloadInput) Payload {
	return Payload{
		Type:   input.Type,
		Data:   input.Data,
		Client: client,
	}
}

type Payload struct {
	Type string
	Data any

	Client *Client
}

func (p *Payload) DataBytes() ([]byte, error) {
	return json.Marshal(p.Data)
}

type ErrorPayload struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type PayloadInput struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type PayloadOutput struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
