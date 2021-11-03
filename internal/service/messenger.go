package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/pkg/logger"
)

type WebSocketMessenger struct {
	rooms  map[string]domain.Room
	logger logger.Logger
}

func NewWebSocketMessenger(rooms map[string]domain.Room, logger logger.Logger) *WebSocketMessenger {
	return &WebSocketMessenger{rooms: rooms, logger: logger}
}

func (m WebSocketMessenger) ConnectClient(client *domain.Client, roomId string) {
	if _, ok := m.rooms[roomId]; ok {
		m.rooms[roomId].Clients[*client] = true
	} else {
		m.createNewRoom(client, roomId)
		go m.messageWriter(roomId, m.rooms[roomId].Broadcast)
	}

	m.logger.Infof("new client: '%s', roomId: '%s'", client.Username, roomId)
	go m.messageReader(client, roomId, m.rooms[roomId].Broadcast)
}

func (m WebSocketMessenger) createNewRoom(client *domain.Client, roomId string) {
	msg := make(chan domain.Message)
	room := domain.Room{
		Id:        roomId,
		Clients:   make(map[domain.Client]bool),
		Broadcast: msg,
	}
	room.Clients[*client] = true
	m.rooms[roomId] = room
}

func (m WebSocketMessenger) messageReader(client *domain.Client, roomId string, messages chan<- domain.Message) {
	defer func() {
		_ = client.Conn.Close()
		delete(m.rooms[roomId].Clients, *client)
		if len(m.rooms[roomId].Clients) == 0 {
			delete(m.rooms, roomId)
			m.logger.Infof("deleted empty room: '%s'", roomId)
		}
		m.logger.Infof("closed connection with client: '%s', roomId: '%s'", client.Username, roomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			m.logger.Debug(err)
			return
		}

		messages <- domain.Message{
			Username: client.Username,
			Text:     string(p),
		}
	}
}

func (m WebSocketMessenger) messageWriter(roomId string, messages <-chan domain.Message) {
	for {
		select {
		case msg := <-messages:
			for client := range m.rooms[roomId].Clients {
				js, err := json.Marshal(msg)
				if err != nil {
					m.logger.Error(err)
					return
				}
				m.logger.Infof("sending message: '%s' to client: '%s'", js, client.Username)

				err = client.Conn.WriteMessage(1, js)
				if err != nil {
					m.logger.Error(err)
					return
				}
			}
		}
	}
}
