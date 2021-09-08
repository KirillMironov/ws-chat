package service

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"log"
)

type Messenger interface {
	ConnectClient(client *domain.Client, roomId string)
	createNewRoom(client *domain.Client, roomId string)
	messageReader(client *domain.Client, roomId string, messages chan<- domain.Message)
	messageWriter(roomId string, messages <-chan domain.Message)
}

type MessengerWS struct {
	rooms map[string]domain.Room
}

func NewMessengerWS(r map[string]domain.Room) *MessengerWS {
	return &MessengerWS{rooms: r}
}

func (m MessengerWS) ConnectClient(client *domain.Client, roomId string) {
	if _, ok := m.rooms[roomId]; ok {
		m.rooms[roomId].Clients[*client] = true
	} else {
		m.createNewRoom(client, roomId)
		go m.messageWriter(roomId, m.rooms[roomId].Broadcast)
	}

	log.Printf("new client: '%s', roomId: '%s'", client.Username, roomId)
	go m.messageReader(client, roomId, m.rooms[roomId].Broadcast)
}

func (m MessengerWS) createNewRoom(client *domain.Client, roomId string) {
	msg := make(chan domain.Message)
	room := domain.Room{
		Id:        roomId,
		Clients:   make(map[domain.Client]bool),
		Broadcast: msg,
	}
	room.Clients[*client] = true
	m.rooms[roomId] = room
}

func (m MessengerWS) messageReader(client *domain.Client, roomId string, messages chan<- domain.Message) {
	defer func() {
		_ = client.Conn.Close()
		delete(m.rooms[roomId].Clients, *client)
		if len(m.rooms[roomId].Clients) == 0 {
			delete(m.rooms, roomId)
			log.Printf("deleted empty room: '%s'", roomId)
		}
		log.Printf("closed connection with client: '%s', roomId: '%s'", client.Username, roomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		messages <- domain.Message{
			Username: client.Username,
			Text:     string(p),
		}
	}
}

func (m MessengerWS) messageWriter(roomId string, messages <-chan domain.Message) {
	for {
		select {
		case msg := <-messages:
			for client := range m.rooms[roomId].Clients {
				js, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("sending message: '%s' to client: '%s'", js, client.Username)

				err = client.Conn.WriteMessage(1, js)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
