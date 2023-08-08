package service

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/repository"
	"chatto/internal/util/containers"
)

type IRoomService interface {
	CreateRoom(room *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error)
	FindRooms() ([]dto.RoomResponse, common.Error)
	FindRoomById(id string) (*dto.RoomResponse, common.Error)
	// FindUserRoomsByUserId Used to get all UserRooms by the user
	FindUserRoomsByUserId(userId string) ([]dto.UserRoomResponse, common.Error)
	// FindRoomsByUserId Used to get all Rooms by the user
	FindRoomsByUserId(userId string) ([]dto.RoomResponse, common.Error)
	// FindRoomMembersById Used to get all users on room
	FindRoomMembersById(roomId string) ([]dto.UserResponse, common.Error)
	// AddUsersInRoom Used to add user into room, This function should check either the room is exists on IRoomService and user on IUserService
	AddUsersInRoom(input dto.UserRoomAddInput, allowPrivate bool) common.Error
	// RemoveUsersInRoom Used to remove user from room, Room should be removed when there are no users left
	RemoveUsersInRoom(input dto.UserRoomRemoveInput) common.Error
	// ClearRoom Used to clear all user from the room, but it will not remove the room
	ClearRoom(roomId string) common.Error
	// ClearUserRooms Used to remove all rooms from the user
	ClearUserRooms(userId string) common.Error
	DeleteRoomById(id string, force bool) common.Error
}

func NewRoomService(roomRepository repository.IRoomRepository, userRoomRepo repository.IUserRoomRepository) IRoomService {
	return roomService{roomRepo: roomRepository, userRoomRepo: userRoomRepo}
}

type roomService struct {
	roomRepo     repository.IRoomRepository
	userRoomRepo repository.IUserRoomRepository
}

func (r roomService) CreateRoom(input *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error) {
	room := dto.NewRoomFromCreateInput(input)
	if room.Private && len(input.MemberIds) != 1 {
		return dto.CreateRoomOutput{}, common.NewError(common.ROOM_CREATE_ERROR, constant.MSG_PRIVATE_ROOM_NOT_2_USER)
	}

	err := r.roomRepo.CreateRoom(room)
	if err != nil {
		return dto.CreateRoomOutput{}, common.NewError(common.ROOM_CREATE_ERROR, constant.MSG_ROOM_CREATION_FAILED)
	}

	// Add all memberIds into room
	userRooms := containers.ConvertSlice(input.MemberIds, func(current *string) model.UserRoom {
		return model.NewUserRoom(room.Id, *current, model.RoomRoleUser)
	})
	err = r.userRoomRepo.AddUsersIntoRoomById(userRooms)
	if err != nil {
		return dto.CreateRoomOutput{}, common.NewError(common.ROOM_CREATE_ERROR, constant.MSG_ROOM_CREATION_FAILED)
	}

	// Get created room
	room, err = r.roomRepo.FindRoomById(room.Id)
	return dto.NewCreateRoomOutput(room), common.NewConditionalError(err, common.ROOM_CREATE_ERROR, constant.MSG_ROOM_CREATION_FAILED)
}

func (r roomService) FindRooms() ([]dto.RoomResponse, common.Error) {
	rooms, err := r.roomRepo.FindRooms()
	roomResponses := containers.ConvertSlice(rooms, dto.NewRoomResponse)
	return roomResponses, common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) FindRoomById(roomId string) (*dto.RoomResponse, common.Error) {
	room, err := r.roomRepo.FindRoomById(roomId)
	roomResponse := dto.NewRoomResponse(room)
	return &roomResponse, common.NewConditionalError(err, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
}

func (r roomService) FindUserRoomsByUserId(userId string) ([]dto.UserRoomResponse, common.Error) {
	userRooms, err := r.userRoomRepo.FindUserRoomsByUserId(userId)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	roomResponse := containers.ConvertSlice(userRooms, dto.NewUserRoomResponse)

	return roomResponse, common.NoError()
}

func (r roomService) FindRoomsByUserId(userId string) ([]dto.RoomResponse, common.Error) {
	rooms, err := r.roomRepo.FindRoomsByUserId(userId)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	roomResponses := containers.ConvertSlice(rooms, dto.NewRoomResponse)
	return roomResponses, common.NoError()
}

func (r roomService) DeleteRoomById(roomId string, force bool) common.Error {
	if !force {
		users, err := r.userRoomRepo.FindUsersByRoomId(roomId)
		if err != nil {
			return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		}

		if !containers.IsEmpty(users) {
			return common.NewError(common.ROOM_NOT_EMPTY_ERROR, constant.MSG_REMOVE_NON_EMPTY_ROOM)
		}
	}

	// TODO: Maybe better using transaction
	err := r.userRoomRepo.RemoveAllUsersFromRoomById(roomId)
	if err != nil {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	err = r.roomRepo.DeleteRoomById(roomId)

	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) FindRoomMembersById(roomId string) ([]dto.UserResponse, common.Error) {
	users, err := r.userRoomRepo.FindUsersByRoomId(roomId)
	if err != nil {
		return nil, common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	userResponse := containers.ConvertSlice(users, dto.NewUserResponse)

	return userResponse, common.NoError()
}

func (r roomService) AddUsersInRoom(input dto.UserRoomAddInput, allowPrivate bool) common.Error {
	// Check room existence
	room, err := r.roomRepo.FindRoomById(input.RoomId)
	if err != nil {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	if !allowPrivate && room.InviteOnly {
		return common.NewError(common.ROOM_IS_PRIVATE_ERROR, constant.MSG_JOIN_PRIVATE_ROOM)
	}

	userRooms := containers.ConvertSlice(input.Users, func(current *dto.UserWithRole) model.UserRoom {
		return model.NewUserRoom(input.RoomId, current.UserId, current.Role)
	})
	err = r.userRoomRepo.AddUsersIntoRoomById(userRooms)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) RemoveUsersInRoom(input dto.UserRoomRemoveInput) common.Error {
	err := r.userRoomRepo.RemoveUsersFromRoomById(input.RoomId, input.UserIds)
	// Remove room when there are no members
	count, err := r.userRoomRepo.GetRoomMemberCountById(input.RoomId)
	if err != nil {
		return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	if count <= 0 {
		err = r.roomRepo.DeleteRoomById(input.RoomId)
		if err != nil {
			return common.NewError(common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
		}
	}

	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) ClearRoom(roomId string) common.Error {
	err := r.userRoomRepo.RemoveAllUsersFromRoomById(roomId)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) ClearUserRooms(userId string) common.Error {
	err := r.userRoomRepo.RemoveAllRoomsFromUserById(userId)
	return common.NewConditionalError(err, common.INTERNAL_SERVER_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}
