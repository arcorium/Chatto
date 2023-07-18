package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(userId string, username string, role Role, conn *websocket.Conn) Client {
	//conn.SetReadLimit(1024)
	return Client{
		Id:              uuid.NewString(),
		UserId:          userId,
		Username:        username,
		Role:            role,
		Conn:            conn,
		Rooms:           make([]string, 1),
		IncomingPayload: make(chan *PayloadOutput, 100), // Make it buffered
	}
}

type Client struct {
	Id       string          `json:"id"` // client id
	UserId   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     Role            `json:"role"`
	Conn     *websocket.Conn `json:"-"`
	Rooms    []string        `json:"groups"` // group_ids

	IncomingPayload chan *PayloadOutput
}

func (c *Client) SendPayload(payload *PayloadOutput) {
	c.IncomingPayload <- payload
}
