package service

import (
	"chatto/internal/constant"
	"chatto/internal/dto"
	"chatto/internal/model"
	"chatto/internal/model/common"
	"chatto/internal/repository"
	"chatto/internal/util/ctrutil"
)

type IRoomService interface {
	CreateRoom(room *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error)
	FindRooms() ([]model.Room, common.Error)
	FindRoomById(id string) (*model.Room, common.Error)
	DeleteRoomById(id string, force bool) common.Error
	AddUserIntoRoom(roomId string, input dto.UserRoomInput) (dto.UserRoomOutput, common.Error)
	RemoveUsersOnRoom(roomId string, input dto.UserRoomInput) common.Error
}

func NewRoomService(roomRepository repository.IRoomRepository, userRoomRepository repository.IUserRoomRepository) IRoomService {
	return roomService{roomRepo: roomRepository, userRoomRepo: userRoomRepository}
}

type roomService struct {
	roomRepo     repository.IRoomRepository
	userRoomRepo repository.IUserRoomRepository
}

func (r roomService) CreateRoom(incomeRoom *dto.CreateRoomInput) (dto.CreateRoomOutput, common.Error) {
	// TODO: Handle memberIds on Input parameter
	room := dto.NewRoomFromCreateInput(incomeRoom)
	err := r.roomRepo.CreateRoom(room)
	if err != nil {
		return dto.CreateRoomOutput{}, common.NewError(common.ROOM_CREATE_ERROR, constant.MSG_ROOM_CREATION_FAILED)
	}

	// Add all memberIds in the room
	userRooms := r.createUserRooms(room.Id, incomeRoom.MemberIds)
	err = r.userRoomRepo.AddUsersIntoRoomById(userRooms)
	if err != nil {
		return dto.CreateRoomOutput{}, common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_CREATION_FAILED)
	}

	// Get the created room
	room, err = r.roomRepo.FindRoomById(room.Id)
	return dto.NewCreateRoomOutput(room), common.NewConditionalError(err, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
}

func (r roomService) FindRooms() ([]model.Room, common.Error) {
	rooms, err := r.roomRepo.FindRooms()
	return rooms, common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) FindRoomById(id string) (*model.Room, common.Error) {
	room, err := r.roomRepo.FindRoomById(id)
	return room, common.NewConditionalError(err, common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
}

func (r roomService) DeleteRoomById(id string, force bool) common.Error {
	if !force {
		users, err := r.userRoomRepo.GetUserIdsOnRoomById(id)
		if err != nil {
			return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
		}

		if !ctrutil.IsEmpty(users) {
			return common.NewError(common.ROOM_NOT_EMPTY_ERROR, constant.MSG_REMOVE_NON_EMPTY_ROOM)
		}
	}

	// TODO: Maybe better using transaction
	err := r.userRoomRepo.RemoveAllUsersFromRoomById(id)
	if err != nil {
		return common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}
	err = r.roomRepo.DeleteRoomById(id)

	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) AddUserIntoRoom(roomId string, input dto.UserRoomInput) (dto.UserRoomOutput, common.Error) {
	// Check room existence
	_, err := r.roomRepo.FindRoomById(roomId)
	if err != nil {
		return dto.UserRoomOutput{}, common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	memberIds, err := r.userRoomRepo.GetUserIdsOnRoomById(roomId)
	if err != nil {
		return dto.UserRoomOutput{}, common.NewError(common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	userRooms := make([]model.UserRoom, 0, len(input.UserIds))
	// Filter room member
	for _, needle := range input.UserIds {
		isMember := false
		for _, haystack := range memberIds {
			if needle == haystack {
				isMember = true
				break
			}
		}
		if !isMember {
			userRooms = append(userRooms, model.NewUserRoom(roomId, needle))
		}
	}

	// Multiple Input
	err = r.userRoomRepo.AddUsersIntoRoomById(userRooms)
	return dto.NewUserRoomOutput(userRooms), common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) RemoveUsersOnRoom(roomId string, input dto.UserRoomInput) common.Error {
	// Check room existence
	_, err := r.roomRepo.FindRoomById(roomId)
	if err != nil {
		return common.NewError(common.ROOM_NOT_FOUND_ERROR, constant.MSG_ROOM_NOT_FOUND)
	}

	// Single Input
	if len(input.UserIds) == 1 {
		err = r.userRoomRepo.RemoveUserFromRoomById(roomId, input.UserIds[0])
		return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
	}

	err = r.userRoomRepo.RemoveUsersFromRoomById(roomId, input.UserIds)
	return common.NewConditionalError(err, common.INTERNAL_REPOSITORY_ERROR, constant.MSG_INTERNAL_SERVER_ERROR)
}

func (r roomService) createUserRooms(roomId string, userIds []string) []model.UserRoom {
	userRooms := make([]model.UserRoom, 0, len(userIds))
	for _, userId := range userIds {
		userRooms = append(userRooms, model.NewUserRoom(roomId, userId))
	}
	return userRooms
}
