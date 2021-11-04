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

func (w WebSocketMessenger) ConnectClient(client *domain.Client) {
	w.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
	go w.messageWriter(client)
	go w.messageReader(client)
}

func (w WebSocketMessenger) messageWriter(client *domain.Client) {
	defer func() {
		_ = client.Conn.Close()
		w.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		message, err := json.Marshal(domain.Message{
			Username: client.Username,
			Text:     string(p),
		})
		if err != nil {
			w.logger.Error(err)
			return
		}

		err = w.repo.Publish(client.RoomId, message)
		if err != nil {
			w.logger.Error(err)
		}
	}
}

func (w WebSocketMessenger) messageReader(client *domain.Client) {
	subscription := w.repo.Subscribe(client.RoomId)

	for {
		select {
		case message := <-subscription:
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				return
			}
			w.logger.Infof("sent message: '%s' to client: '%s'", message.Payload, client.Username)
		}
	}
}
