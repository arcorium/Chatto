package model

import (
	"encoding/json"
)

const (
	PayloadTyping           = "typing"
	PayloadMessage          = "chat"
	PayloadNotification     = "notif"
	PayloadCreateRoom       = "create-room"
	PayloadJoinRoom         = "join-room"
	PayloadLeaveRoom        = "leave-room"
	PayloadInviteToRoom     = "invite-room"
	PayloadKickFromRoom     = "kick-room"
	PayloadGetUsers         = "get-users"
	PayloadGetChats         = "get-chats"
	PayloadGetNotifications = "get-notifs"
	PayloadGetUserRooms     = "user-rooms"
	PayloadErrorResponse    = "error"
	PayloadSuccessResponse  = "success"
)

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
	output := PayloadOutput{
		Type: types,
	}

	if data != nil {
		output.Data = data
	}
	return output
}

type Payload struct {
	Type string
	Data any

	Sender *Client
}

type PayloadInput struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func PayloadData[T any](payload *Payload) (T, error) {
	var t T
	bytes, err := json.Marshal(payload.Data)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(bytes, &t)
	return t, err
}

type ErrorPayload struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type PayloadOutput struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}
