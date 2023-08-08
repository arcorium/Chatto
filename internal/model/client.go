package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewClient(userId string, username string, role Role, conn *websocket.Conn) Client {
	client := Client{
		Id:              uuid.NewString(),
		UserId:          userId,
		Username:        username,
		Role:            role,
		Conn:            conn,
		IncomingPayload: make(chan *PayloadOutput, 100), // Make it buffered
	}

	return client
}

type Client struct {
	Id       string          `json:"id"`
	UserId   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     Role            `json:"role"`
	Conn     *websocket.Conn `json:"-"`

	IncomingPayload chan *PayloadOutput `json:"-"`
}

func (c *Client) SendPayload(payload *PayloadOutput) {
	c.IncomingPayload <- payload
}
