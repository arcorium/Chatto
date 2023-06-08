package manager

import (
	"errors"
	"strings"
	"sync"

	"chatto/internal/model"
)

type ClientList map[string]*model.Client

func NewClientManager() ClientManager {
	return ClientManager{
		clientsMutex: sync.RWMutex{},
		Clients:      make(ClientList),
	}
}

type ClientManager struct {
	clientsMutex sync.RWMutex
	Clients      ClientList
}

func (m *ClientManager) AddClients(clients ...*model.Client) {
	m.clientsMutex.Lock()
	defer m.clientsMutex.Unlock()
	for _, c := range clients {
		m.Clients[c.Id] = c
	}
}

func (m *ClientManager) GetClientsByUserId(userId string) []*model.Client {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()
	clients := make([]*model.Client, 0, 10)
	for _, c := range m.Clients {
		if c.UserId == userId {
			clients = append(clients, c)
		}
	}
	return clients
}

func (m *ClientManager) GetClientById(clientId string) (*model.Client, error) {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()
	client, ok := m.Clients[clientId]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (m *ClientManager) GetClientsByUsername(username string) []*model.Client {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	username = strings.ToLower(username)
	clientMaps := make(map[string]*model.Client)
	for _, c := range m.Clients {
		if strings.Contains(strings.ToLower(c.Username), username) {
			clientMaps[c.UserId] = c
		}
	}

	// get the values
	clients := make([]*model.Client, 0, len(clientMaps))
	for _, v := range clientMaps {
		clients = append(clients, v)
	}

	return clients
}

func (m *ClientManager) RemoveClientsByUserId(userId string) {
	clients := m.GetClientsByUserId(userId)
	if len(clients) == 0 {
		return
	}

	m.clientsMutex.Lock()
	for _, c := range clients {
		delete(m.Clients, c.Id)
	}
	m.clientsMutex.Unlock()
}

func (m *ClientManager) RemoveClientById(clientId string) {
	// TODO: Maybe getting the client is not necessary and just delete the clientId
	client, err := m.GetClientById(clientId)
	if err != nil {
		return
	}
	m.clientsMutex.Lock()
	delete(m.Clients, client.Id)
	m.clientsMutex.Unlock()
}
