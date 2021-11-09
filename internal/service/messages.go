package service

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/repository"
	"github.com/KirillMironov/ws-chat/pkg/logger"
	"github.com/gorilla/websocket"
	"sync"
)

type MessagesService struct {
	clientsRepo  repository.Clients
	messagesRepo repository.Messages
	logger       logger.Logger
}

func NewMessagesService(clientsRepo repository.Clients, messagesRepo repository.Messages,
	logger logger.Logger) *MessagesService {
	return &MessagesService{clientsRepo: clientsRepo, messagesRepo: messagesRepo, logger: logger}
}

func (m MessagesService) MessageWriter(client *domain.Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		var message = domain.Message{Event: domain.ChatMessage}
		message.Payload.Username = client.Username
		message.Payload.Text = string(p)

		js, err := json.Marshal(message)
		if err != nil {
			m.logger.Error(err)
			return
		}

		err = m.messagesRepo.PublishMessage(client.RoomId, js)
		if err != nil {
			m.logger.Error(err)
			return
		}
	}
}

func (m MessagesService) MessageReader(client *domain.Client, done <-chan struct{}, wg *sync.WaitGroup) {
	messagesSubscription := m.messagesRepo.SubscribeToMessages(client.RoomId)
	defer messagesSubscription.Unsubscribe(context.Background())
	activeClientsSubscription := m.clientsRepo.SubscribeToActiveClients(client.RoomId)
	defer activeClientsSubscription.Unsubscribe(context.Background())
	wg.Done()

	for {
		select {
		case message := <-messagesSubscription.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				return
			}
			m.logger.Infof("sent message: '%s' to client: '%s'", message.Payload, client.Username)
		case message := <-activeClientsSubscription.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				m.logger.Error(err)
				return
			}
			m.logger.Infof("sent message: '%s' to client: '%s'", message.Payload, client.Username)
		case <-done:
			return
		}
	}
}
