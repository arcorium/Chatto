package service

import "chatto/internal/repository"

func NewRoomService(roomRepository repository.IRoomRepository) IRoomService {
	return roomService{roomRepo: roomRepository}
}

type roomService struct {
	roomRepo repository.IRoomRepository
}

func (r roomService) CreateRoom() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) FindRooms() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) FindRoomByName() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) FindRoomById() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) DeleteRoom() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) AddUserIntoRoom() {
	//TODO implement me
	panic("implement me")
}

func (r roomService) DeleteUserFromRoom() {
	//TODO implement me
	panic("implement me")
}
