package model

import (
	"time"

	"github.com/google/uuid"
)

func NewRoomPayload(room *Room) *Payload {
	return &Payload{
		Type:     PayloadCreateRoom,
		ClientId: "server",
		Data: RoomResponse{
			Id:      room.Id,
			Name:    room.Name,
			Private: room.Private,
		},
	}
}

func NewRoom(name string, private bool, clients ...*Client) Room {
	room := Room{
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
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Private   bool               `json:"private"` // Used for only invite
	CreatedAt time.Time          `json:"created_at"`
	Clients   map[string]*Client `json:"-"`
}

type RoomResponse struct {
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

func (r *Room) BroadcastPayloadExceptUserId(payload *Payload, userId string) {
	for _, c := range r.Clients {
		if c.UserId == userId {
			continue
		}
		c.SendPayload(payload)
	}
}

func (r *Room) BroadcastPayloadExceptClientId(payload *Payload, clientId string) {
	for _, c := range r.Clients {
		if c.Id == clientId {
			continue
		}
		c.SendPayload(payload)
	}
}

// PrivateChat Used for user to user chat, server will create the room and set the private into true
type PrivateChat struct {
	Opponent string `json:"opponent"`
}

// CreateRoom Used to create new room with members, it is needed due to no implicit feature to create the room when trying to join unlisted room
type CreateRoom struct {
	Name    string   `json:"name"`
	Private bool     `json:"private"`
	Members []string `json:"members,omitempty"` // Initial members. TODO: Members should be on sender friends
}

// JoinRoom Used to join room, room by the roomId should check the private
type JoinRoom struct {
	RoomId string `json:"room_id"`
}

// InviteRoom Used to invite another clients to join the room, which will send the client either to accept or not (For now all the invited clients always accepting)
type InviteRoom struct {
	RoomId   string   `json:"room_id"`
	Receiver []string `json:"receiver"`
}
