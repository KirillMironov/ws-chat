package service

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
)

type ClientsService struct {
	repository      ports.ClientsRepository
	messagesService ports.MessagesService
	logger          ports.Logger
}

func NewClientsService(repository ports.ClientsRepository, messagesService ports.MessagesService,
	logger ports.Logger) *ClientsService {
	return &ClientsService{repository: repository, messagesService: messagesService, logger: logger}
}

func (c ClientsService) Connect(client *domain.Client) {
	done := make(chan struct{})
	go c.messagesService.Writer(client, done)
	defer c.Disconnect(client, done)

	<-done

	err := c.repository.Add(client)
	if err != nil {
		c.logger.Error(err)
		return
	}

	c.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
	c.messagesService.Reader(client)
}

func (c ClientsService) Disconnect(client *domain.Client, done chan<- struct{}) {
	done <- struct{}{}
	client.Conn.Close()
	err := c.repository.Remove(client)
	if err != nil {
		c.logger.Error(err)
		return
	}
	c.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
}
