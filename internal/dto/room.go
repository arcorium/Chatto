package dto

import (
	"time"

	"chatto/internal/model"
	"chatto/internal/util/containers"
	"github.com/google/uuid"
)

func NewCreateRoomOutput(room *model.Room) CreateRoomOutput {
	return CreateRoomOutput{
		Id:         room.Id,
		Name:       room.Name,
		InviteOnly: room.InviteOnly,
		Private:    room.Private,
	}
}

func NewRoomFromCreateInput(room *CreateRoomInput) *model.Room {
	return &model.Room{
		Id:          uuid.NewString(),
		Name:        room.Name,
		Description: room.Description,
		InviteOnly:  room.InviteOnly,
		Private:     room.Private,
		CreatedAt:   time.Now(),
	}
}

// CreateRoomInput Used to create new room with members, it is needed due to no implicit feature to create the room when trying to join unlisted room
type CreateRoomInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"desc"`
	InviteOnly  bool     `json:"invite_only"`
	Private     bool     `json:"private"`
	MemberIds   []string `json:"member_ids"` // Initial members. TODO: MemberIds should be on sender friends
}

// RoomInput Used to join and leave room, room by the roomId should check the private
type RoomInput struct {
	RoomId string `json:"room_id"`
}

// MemberRoomInput Used to invite another clients to join the room, which will send the client either to accept or not (For now all the invited clients always accepting)
type MemberRoomInput struct {
	RoomId  string   `json:"room_id"`
	UserIds []string `json:"user_ids"`
}

func NewChatRoomFromOutput(roomOutput *CreateRoomOutput, members ...*model.Client) *model.ChatRoom {
	var room model.ChatRoom
	if roomOutput.Private {
		room = model.NewPrivateChatRoom(roomOutput.Id, roomOutput.Name, "", roomOutput.InviteOnly)
	} else {
		room = model.NewChatRoom(roomOutput.Id, roomOutput.Name, "", roomOutput.InviteOnly)
	}

	room.AddClientsWithSameRole(model.RoomRoleUser, members...)
	return &room
}

type CreateRoomOutput struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	InviteOnly bool   `json:"invite_only"`
	Private    bool   `json:"private"`
}

func NewUserRoomOutput(userResponses []model.UserRoom) UserRoomOutput {
	users := containers.ConvertSlice(userResponses, func(current *model.UserRoom) UserWithRole {
		return UserWithRole{
			UserId: current.UserId,
			Role:   current.UserRole,
		}
	})
	return UserRoomOutput{Users: users}
}

type UserRoomOutput struct {
	Users []UserWithRole `json:"users"`
}

func NewRoomResponse(room *model.Room) RoomResponse {
	return RoomResponse{
		Id:          room.Id,
		Name:        room.Name,
		Description: room.Description,
		InviteOnly:  room.InviteOnly,
		Private:     room.Private,
	}
}

type RoomResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	InviteOnly  bool   `json:"invite_only"`
	Private     bool   `json:"private"`
}
