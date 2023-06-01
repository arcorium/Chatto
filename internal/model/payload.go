package model

import (
	"encoding/json"
)

const (
	PayloadMessage      = "message"
	PayloadNotification = "notification"
	PayloadStartChat    = "start-chat"
	PayloadCreateRoom   = "create-room"
	PayloadJoinRoom     = "join-room"
	PayloadLeaveRoom    = "leave-room"
	PayloadError        = "error"
)

func Decode[T any](bytes []byte) (T, error) {
	var t T
	err := json.Unmarshal(bytes, &t)
	return t, err
}

func NewErrorResponsePayload(message string) Payload {
	return Payload{
		Type: PayloadError,
		Data: struct {
			Message string `json:"message"`
		}{message},
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
