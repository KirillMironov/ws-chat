package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/repository"
	"github.com/KirillMironov/ws-chat/pkg/logger"
	"sync"
)

type ClientsService struct {
	repo            repository.Clients
	messagesService Messages
	logger          logger.Logger
}

func NewClientsService(repo repository.Clients, messagesService Messages, logger logger.Logger) *ClientsService {
	return &ClientsService{repo: repo, messagesService: messagesService, logger: logger}
}

func (c ClientsService) ConnectClient(client *domain.Client) {
	done := make(chan struct{})
	defer c.disconnectClient(client, done)

	err := c.repo.AddActiveClient(client.RoomId, client.Username)
	if err != nil {
		c.logger.Error(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go c.messagesService.MessageReader(client, done, &wg)
	wg.Wait()

	c.updateActiveClients(client)
	c.messagesService.MessageWriter(client)

	c.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) disconnectClient(client *domain.Client, done chan<- struct{}) {
	client.Conn.Close()
	done <- struct{}{}
	err := c.repo.RemoveActiveClient(client.RoomId, client.Username)
	if err != nil {
		c.logger.Error(err)
	}
	c.updateActiveClients(client)
	c.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) updateActiveClients(client *domain.Client) {
	clients, err := c.repo.GetActiveClients(client.RoomId)
	if err != nil {
		c.logger.Error(err)
		return
	}

	var message = domain.Message{Event: domain.ActiveClients}
	message.Payload.Text = clients

	js, err := json.Marshal(message)
	if err != nil {
		c.logger.Error(err)
		return
	}

	err = c.repo.PublishActiveClients(client.RoomId, js)
	if err != nil {
		c.logger.Error(err)
		return
	}
}
