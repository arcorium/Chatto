package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(userId string, username string, role Role, status ClientStatus, conn *websocket.Conn) Client {
	//conn.SetReadLimit(1024)
	return Client{
		Id:              uuid.NewString(),
		UserId:          userId,
		Username:        username,
		Role:            role,
		Conn:            conn,
		Status:          status,
		Rooms:           make([]string, 1),
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
	Id       string          `json:"id"` // client id
	UserId   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     Role            `json:"role"`
	Conn     *websocket.Conn `json:"-"`
	Status   ClientStatus    `json:"status"`
	Rooms    []string        `json:"groups"` // group_ids

	IncomingPayload chan *Payload
}

func (c *Client) SendPayload(payload *Payload) {
	c.IncomingPayload <- payload
}
