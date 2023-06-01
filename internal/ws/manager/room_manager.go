package manager

import (
	"errors"
	"sync"

	"server_client_chat/internal/model"
)

type RoomList map[string]*model.Room

func NewRoomManager() RoomManager {
	return RoomManager{
		roomsMutex: sync.RWMutex{},
		rooms:      make(RoomList),
	}
}

type RoomManager struct {
	roomsMutex sync.RWMutex
	rooms      RoomList
}

func (r *RoomManager) AddRooms(rooms ...*model.Room) {
	r.roomsMutex.Lock()
	for _, room := range rooms {
		r.rooms[room.Id] = room
	}
	r.roomsMutex.Unlock()
}

func (r *RoomManager) GetRoomById(roomId string) (*model.Room, error) {
	r.roomsMutex.RLock()
	defer r.roomsMutex.RUnlock()
	room, ok := r.rooms[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (r *RoomManager) GetRoomByName(name string) (*model.Room, error) {
	r.roomsMutex.RLock()
	defer r.roomsMutex.RUnlock()
	for _, room := range r.rooms {
		if room.Name == name {
			return room, nil
		}
	}
	return nil, errors.New("room not found")
}

func (r *RoomManager) RemoveRoomById(roomId string) error {
	room, err := r.GetRoomById(roomId)
	if err != nil {
		return err
	}

	r.roomsMutex.Lock()
	delete(r.rooms, room.Id)
	r.roomsMutex.Unlock()
	return nil
}

func (r *RoomManager) RemoveRoomByName(name string) error {
	room, err := r.GetRoomByName(name)
	if err != nil {
		return err
	}

	r.roomsMutex.Lock()
	delete(r.rooms, room.Id)
	r.roomsMutex.Unlock()
	return nil
}
