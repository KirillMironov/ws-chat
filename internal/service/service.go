package service

import "github.com/KirillMironov/ws-chat/domain"

type Messenger interface {
	ConnectClient(client *domain.Client, roomId string)
	createNewRoom(client *domain.Client, roomId string)
	messageReader(client *domain.Client, roomId string, messages chan<- domain.Message)
	messageWriter(roomId string, messages <-chan domain.Message)
}
