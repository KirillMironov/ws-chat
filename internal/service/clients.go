package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
	"sync"
)

type ClientsService struct {
	repo            ports.ClientsRepository
	messagesService ports.MessagesService
	logger          ports.Logger
}

func NewClientsService(repo ports.ClientsRepository, messagesService ports.MessagesService,
	logger ports.Logger) *ClientsService {
	return &ClientsService{repo: repo, messagesService: messagesService, logger: logger}
}

func (c ClientsService) ConnectClient(client *domain.Client) {
	done := make(chan struct{})
	defer c.DisconnectClient(client, done)

	err := c.repo.AddActiveClient(client.RoomId, client.Username)
	if err != nil {
		c.logger.Error(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go c.messagesService.MessageReader(client, done, &wg)
	wg.Wait()

	c.UpdateActiveClients(client)
	c.messagesService.MessageWriter(client)

	c.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) DisconnectClient(client *domain.Client, done chan<- struct{}) {
	client.Conn.Close()
	done <- struct{}{}
	err := c.repo.RemoveActiveClient(client.RoomId, client.Username)
	if err != nil {
		c.logger.Error(err)
	}
	c.UpdateActiveClients(client)
	c.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) UpdateActiveClients(client *domain.Client) {
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
