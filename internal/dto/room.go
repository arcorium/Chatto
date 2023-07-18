package dto

import (
	"chatto/internal/model"
	"chatto/internal/util/ctrutil"
	"github.com/google/uuid"
	"time"
)

func NewCreateRoomOutput(room *model.Room) CreateRoomOutput {
	return CreateRoomOutput{
		Id:      room.Id,
		Name:    room.Name,
		Private: room.Private,
	}
}

func NewRoomFromCreateInput(room *CreateRoomInput) *model.Room {
	// TODO: Handle userids
	return &model.Room{
		Id:        uuid.NewString(),
		Name:      room.Name,
		Private:   room.Private,
		CreatedAt: time.Now(),
	}
}

// CreateRoomInput Used to create new room with members, it is needed due to no implicit feature to create the room when trying to join unlisted room
type CreateRoomInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"desc"`
	Private     bool     `json:"private"`
	MemberIds   []string `json:"members"` // Initial members. TODO: MemberIds should be on sender friends
}

// JoinRoomInput Used to join room, room by the roomId should check the private
type JoinRoomInput struct {
	RoomId string `json:"room_id"`
}

type LeaveRoomInput struct {
	RoomId string `json:"room_id"`
}

// InviteRoomInput Used to invite another clients to join the room, which will send the client either to accept or not (For now all the invited clients always accepting)
type InviteRoomInput struct {
	RoomId  string   `json:"room_id"`
	UserIds []string `json:"user_ids"`
}

func NewChatRoom(roomOutput *CreateRoomOutput, clients ...*model.Client) *model.ChatRoom {
	room := &model.ChatRoom{
		Id:      roomOutput.Id,
		Name:    roomOutput.Name,
		Private: roomOutput.Private,
		Clients: make(map[string]*model.Client),
	}
	room.AddClients(clients...)
	return room
}

type CreateRoomOutput struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type UserRoomInput struct {
	UserIds []string `json:"users"`
}

func NewUserRoomOutput(userResponses []model.UserRoom) UserRoomOutput {
	userIds := ctrutil.ConvertSliceType(userResponses, func(current *model.UserRoom) string {
		return current.UserId
	})
	return UserRoomOutput{UserIds: userIds}
}

type UserRoomOutput struct {
	UserIds []string `json:"users"`
}
