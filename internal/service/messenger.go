package service

import (
	"context"
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
	done := make(chan struct{})

	err := w.repo.AddActiveUser(client.RoomId, client.Username)
	if err != nil {
		w.logger.Error(err)
		return
	}

	go w.messageWriter(client, done)
	go w.messageReader(client, done)

	w.logger.Infof("new client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (w WebSocketMessenger) disconnectClient(client *domain.Client, done chan<- struct{}) {
	client.Conn.Close()
	err := w.repo.RemoveActiveUser(client.RoomId, client.Username)
	if err != nil {
		w.logger.Error(err)
	}
	w.publishActiveUsers(client)
	done <- struct{}{}
	w.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, client.RoomId)
}

func (w WebSocketMessenger) publishActiveUsers(client *domain.Client) {
	users, err := w.repo.GetActiveUsers(client.RoomId)
	if err != nil {
		w.logger.Error(err)
		return
	}

	var message = domain.Message{Event: domain.ActiveClientsCounter}
	message.Payload.Text = users

	js, err := json.Marshal(message)
	if err != nil {
		w.logger.Error(err)
		return
	}

	err = w.repo.PublishActiveUsers(client.RoomId, js)
	if err != nil {
		w.logger.Error(err)
		return
	}
}

func (w WebSocketMessenger) messageWriter(client *domain.Client, done chan<- struct{}) {
	defer w.disconnectClient(client, done)

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
			w.logger.Error(err)
			return
		}

		err = w.repo.PublishMessage(client.RoomId, js)
		if err != nil {
			w.logger.Error(err)
			return
		}
	}
}

func (w WebSocketMessenger) messageReader(client *domain.Client, done <-chan struct{}) {
	messagesSubscription := w.repo.SubscribeToMessages(client.RoomId)
	defer messagesSubscription.Unsubscribe(context.Background(), client.RoomId)
	activeUsersSubscription := w.repo.SubscribeToActiveUsers(client.RoomId)
	defer activeUsersSubscription.Unsubscribe(context.Background(), client.RoomId)
	w.publishActiveUsers(client)

	for {
		select {
		case message := <-messagesSubscription.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				return
			}
			w.logger.Infof("sent message: '%s' to client: '%s'", message.Payload, client.Username)
		case message := <-activeUsersSubscription.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			if err != nil {
				w.logger.Error(err)
				return
			}
		case <-done:
			return
		}
	}
}
