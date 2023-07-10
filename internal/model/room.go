package model

import (
	"time"

	"chatto/internal/util"
	"github.com/google/uuid"
)

func NewOutcomeCreateRoom(room *Room) OutcomeCreateRoom {
	return OutcomeCreateRoom{
		Id:      room.Id,
		Name:    room.Name,
		Private: room.Private,
	}
}

func NewRoom(name string, private bool, clients ...*Client) *Room {
	room := &Room{
		Id:        uuid.NewString(),
		Name:      name,
		Private:   private,
		CreatedAt: time.Now(),
		Clients:   make(map[string]*Client),
	}
	room.AddClients(clients...)
	return room
}

type Room struct {
	Id        string             `json:"id" gorm:"primaryKey;type:uuid;not null"`
	Name      string             `json:"name" gorm:"not null"`
	Private   bool               `json:"private" gorm:"not null"` // Used for invite only
	Users     []string           `json:"users" gorm:"type:uuid[]"`
	CreatedAt time.Time          `json:"created_at"`
	Clients   map[string]*Client `json:"-" gorm:"-"`
}

type OutcomeCreateRoom struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (r *Room) IsExist(client *Client) bool {
	return r.Clients[client.UserId] != nil
}

func (r *Room) AddClients(clients ...*Client) {
	for _, c := range clients {
		r.Clients[c.Id] = c
	}
}

func (r *Room) RemoveClients(clients ...*Client) {
	for _, c := range clients {
		delete(r.Clients, c.Id)
	}
}

func (r *Room) RemoveClientsByUserId(userId string) {
	for _, c := range r.Clients {
		if userId == c.UserId {
			delete(r.Clients, c.Id)
		}
	}
}

func (r *Room) BroadcastPayload(payload *Payload) {
	for _, c := range r.Clients {
		c.SendPayload(payload)
	}
}

func (r *Room) Broadcast(payload *Payload, excludeIds ...string) {
	for _, client := range r.Clients {
		// Check if the client id is excluded
		if !util.IsExist(excludeIds, func(current string) bool {
			return current == client.Id
		}) {
			client.SendPayload(payload)
		}
	}
}

// IncomeCreateRoom Used to create new room with members, it is needed due to no implicit feature to create the room when trying to join unlisted room
type IncomeCreateRoom struct {
	Name      string   `json:"name"`
	Private   bool     `json:"private"`
	MemberIds []string `json:"members,omitempty"` // Initial members. TODO: MemberIds should be on sender friends
}

// IncomeJoinRoom Used to join room, room by the roomId should check the private
type IncomeJoinRoom struct {
	RoomId string `json:"room_id"`
}

type IncomeLeaveRoom struct {
	RoomId string `json:"room_id"`
}

// IncomeInviteRoom Used to invite another clients to join the room, which will send the client either to accept or not (For now all the invited clients always accepting)
type IncomeInviteRoom struct {
	RoomId  string   `json:"room_id"`
	UserIds []string `json:"user_ids"`
}
