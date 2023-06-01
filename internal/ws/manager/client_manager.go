package manager

import (
	"errors"
	"sync"

	"server_client_chat/internal/model"
)

type ClientList map[string]*model.Client

func NewClientManager() ClientManager {
	return ClientManager{
		clientsMutex: sync.RWMutex{},
		clients:      make(ClientList),
	}
}

type ClientManager struct {
	clientsMutex sync.RWMutex
	clients      ClientList
}

func (m *ClientManager) AddClients(clients ...*model.Client) {
	m.clientsMutex.Lock()
	defer m.clientsMutex.Unlock()
	for _, c := range clients {
		m.clients[c.Id] = c
	}
}

func (m *ClientManager) GetClientsByUserId(userId string) []*model.Client {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()
	clients := make([]*model.Client, 0, 10)
	for _, c := range m.clients {
		if c.UserId == userId {
			clients = append(clients, c)
		}
	}
	return clients
}

func (m *ClientManager) GetClientById(clientId string) (*model.Client, error) {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()
	client, ok := m.clients[clientId]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (m *ClientManager) RemoveClientByUserId(userId string) {
	clients := m.GetClientsByUserId(userId)
	if len(clients) == 0 {
		return
	}

	m.clientsMutex.Lock()
	for _, c := range clients {
		delete(m.clients, c.Id)
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
	delete(m.clients, client.Id)
	m.clientsMutex.Unlock()
}
