package model

import (
	"time"

	"chatto/internal/util/containers"
)

type RoomRole string

const (
	RoomRoleUser  RoomRole = "user"
	RoomRoleAdmin RoomRole = "admin"
)

type Room struct {
	Id          string `gorm:"primaryKey;type:uuid;not null"`
	Name        string `gorm:"not null"`
	Description string
	InviteOnly  bool `gorm:"not null"` // Used for invite only
	Private     bool `gorm:"not null"` // Used for eiter this room is private chat (2 users) or not

	CreatedAt time.Time
}

func newChatRoom(id, name, desc string, inviteOnly, private bool) ChatRoom {
	return ChatRoom{
		Id:          id,
		Name:        name,
		Description: desc,
		InviteOnly:  inviteOnly,
		Private:     private,
		clients:     make(map[string]*Client, 0),
		roles:       make(map[string]RoomRole, 0),
	}
}

func NewChatRoom(id, name, desc string, inviteOnly bool) ChatRoom {
	return newChatRoom(id, name, desc, inviteOnly, false)
}

func NewPrivateChatRoom(id, name, desc string, inviteOnly bool) ChatRoom {
	return newChatRoom(id, name, desc, inviteOnly, true)
}

type ChatRoom struct {
	Id          string
	Name        string
	Description string
	InviteOnly  bool
	Private     bool
	clients     map[string]*Client  // key : clientId
	roles       map[string]RoomRole // key : userId
}

// IsClientExist Used to check single client if it is already on the room
func (r *ChatRoom) IsClientExist(client *Client) bool {
	return r.clients[client.Id] != nil
}

func (r *ChatRoom) GetRoleByUserId(userId string) (RoomRole, bool) {
	role, exist := r.roles[userId]
	return role, exist
}

// IsUserExist Used to check if the user is already on the room
func (r *ChatRoom) IsUserExist(userId string) bool {
	return containers.MapIsExist(r.clients, func(key string, val *Client) bool {
		return val.UserId == userId
	})
}

func (r *ChatRoom) Clients() []*Client {
	return containers.MapValues(r.clients)
}

func (r *ChatRoom) UserIds() []string {
	return containers.MapKeys(r.roles)
}

func (r *ChatRoom) AddClient(client *Client, role RoomRole) {
	r.clients[client.Id] = client
	_, exist := r.roles[client.UserId]
	if !exist {
		r.roles[client.UserId] = role
	}
}

func (r *ChatRoom) AddClientsWithSameRole(role RoomRole, clients ...*Client) {
	for _, c := range clients {
		r.AddClient(c, role)
	}
}

func (r *ChatRoom) RemoveClient(client *Client) {
	delete(r.clients, client.Id)
	if containers.IsEmpty(r.GetClientsByUserId(client.UserId)) {
		delete(r.roles, client.UserId)
	}
}

func (r *ChatRoom) GetClientsByUserId(userId string) []*Client {
	clients := make([]*Client, 0, 3)
	for _, c := range r.clients {
		if userId == c.UserId {
			clients = append(clients, c)
		}
	}
	return clients
}

func (r *ChatRoom) RemoveClientsByUserId(userId string) {
	for _, c := range r.clients {
		if userId == c.UserId {
			delete(r.clients, c.Id)
		}
	}
	delete(r.roles, userId)
}

// Broadcast Used for send payload to all users in room except clients from parameter excludeClientIds
func (r *ChatRoom) Broadcast(payload *PayloadOutput, excludeClientIds ...string) {
	for _, client := range r.clients {
		// Check if the client id is excluded
		if !containers.IsExist(excludeClientIds, func(current *string) bool {
			return *current == client.Id
		}) {
			client.SendPayload(payload)
		}
	}
}
