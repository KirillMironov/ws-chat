package service

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
	"github.com/gorilla/websocket"
)

type MessagesService struct {
	messagesRepo ports.MessagesRepository
	logger       ports.Logger
}

func NewMessagesService(messagesRepo ports.MessagesRepository, logger ports.Logger) *MessagesService {
	return &MessagesService{messagesRepo: messagesRepo, logger: logger}
}

func (m MessagesService) Reader(client *domain.Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		err = m.messagesRepo.Publish(client, p)
		if err != nil {
			m.logger.Error(err)
			return
		}
	}
}

func (m MessagesService) Writer(client *domain.Client, done chan struct{}) {
	sub := m.messagesRepo.Subscribe(client.RoomId)
	defer sub.Close()
	done <- struct{}{}

	for {
		select {
		case message := <-sub.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				return
			}
			m.logger.Infof("sent message: '%s' to client: '%s'", message.Payload, client.Username)
		case <-done:
			return
		}
	}
}
