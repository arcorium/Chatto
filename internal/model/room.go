package model

import (
	"chatto/internal/util/ctrutil"
	"time"
)

type Room struct {
	Id          string `gorm:"primaryKey;type:uuid;not null"`
	Name        string `gorm:"not null"`
	Description string
	Private     bool `gorm:"not null"` // Used for invite only

	CreatedAt time.Time
}

type ChatRoom struct {
	Id          string
	Name        string
	Description string
	Private     bool
	Clients     map[string]*Client
}

// IsClientExist Used to check single client if it is already on the room
func (r *ChatRoom) IsClientExist(client *Client) bool {
	return r.Clients[client.Id] != nil
}

// IsUserExist Used to check if the user is already on the room
func (r *ChatRoom) IsUserExist(userId string) bool {
	return ctrutil.MapIsExist(r.Clients, func(key string, val *Client) bool {
		return val.UserId == userId
	})
}

func (r *ChatRoom) AddClients(clients ...*Client) {
	for _, c := range clients {
		r.Clients[c.Id] = c
	}
}

func (r *ChatRoom) RemoveClients(clients ...*Client) {
	for _, c := range clients {
		delete(r.Clients, c.Id)
	}
}

func (r *ChatRoom) RemoveClientsByUserId(userId string) {
	for _, c := range r.Clients {
		if userId == c.UserId {
			delete(r.Clients, c.Id)
		}
	}
}

func (r *ChatRoom) Broadcast(payload *PayloadOutput, excludeIds ...string) {
	for _, client := range r.Clients {
		// Check if the client id is excluded
		if !ctrutil.IsExist(excludeIds, func(current *string) bool {
			return *current == client.Id
		}) {
			client.SendPayload(payload)
		}
	}
}
