package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(userId string, status ClientStatus, conn *websocket.Conn) Client {
	//conn.SetReadLimit(1024)
	return Client{
		Id:              uuid.NewString(),
		UserId:          userId,
		Conn:            conn,
		Status:          status,
		IncomingPayload: make(chan *Payload, 100), // Make it buffered
	}
}

type ClientStatus string

const (
	ClientStatusOnline     ClientStatus = "on"
	ClientStatusOffline                 = "off"
	ClientStatusRegister                = "reg"
	ClientStatusUnregister              = "unreg"
)

type Client struct {
	Id       string          `json:"id"`
	UserId   string          `json:"user_id"`
	Username string          `json:"username"`
	Conn     *websocket.Conn `json:"-"`
	Status   ClientStatus    `json:"status"`
	Rooms    []string        `json:"groups"` // group_ids

	IncomingPayload chan *Payload `json:"-"`
}

func (c *Client) SendPayload(payload *Payload) {
	c.IncomingPayload <- payload
}
