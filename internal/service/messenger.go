package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/repository"
	"github.com/KirillMironov/ws-chat/pkg/logger"
	"github.com/gorilla/websocket"
)

type WebSocketMessenger struct {
	repo   repository.Messages
	logger logger.Logger
}

func NewWebSocketMessenger(repo repository.Messages, logger logger.Logger) *WebSocketMessenger {
	return &WebSocketMessenger{repo: repo, logger: logger}
}

func (m WebSocketMessenger) ConnectClient(client *domain.Client) {
	m.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
	go m.messageWriter(client)
	go m.messageReader(client)
}

func (m WebSocketMessenger) messageWriter(client *domain.Client) {
	defer func() {
		_ = client.Conn.Close()
		m.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		message, err := json.Marshal(map[string]string{client.Username: string(p)})
		if err != nil {
			m.logger.Error(err)
			return
		}

		err = m.repo.SendMessage(client.RoomId, message)
		if err != nil {
			m.logger.Error(err)
		}
	}
}

func (m WebSocketMessenger) messageReader(client *domain.Client) {
	subscription := m.repo.GetMessages(client.RoomId)

	for {
		select {
		case message := <-subscription:
			m.logger.Infof("sending message: '%s' to client: '%s'", message.Payload, client.Username)
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				m.logger.Error(err)
			}
		}
	}
}
