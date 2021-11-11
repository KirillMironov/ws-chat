package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
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

func (c ClientsService) Connect(client *domain.Client) {
	done := make(chan struct{})
	defer c.Disconnect(client, done)

	err := c.repo.Add(client)
	if err != nil {
		c.logger.Error(err)
		return
	}

	go c.messagesService.Writer(client, done)
	<-done

	c.UpdateConnected(client)
	c.messagesService.Reader(client)

	c.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) Disconnect(client *domain.Client, done chan<- struct{}) {
	client.Conn.Close()
	done <- struct{}{}
	err := c.repo.Remove(client)
	if err != nil {
		c.logger.Error(err)
	}
	c.UpdateConnected(client)
	c.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (c ClientsService) UpdateConnected(client *domain.Client) {
	clients, err := c.repo.GetConnected(client)
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

	err = c.repo.Publish(client, js)
	if err != nil {
		c.logger.Error(err)
		return
	}
}
