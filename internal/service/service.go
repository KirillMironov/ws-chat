package service

import "github.com/KirillMironov/ws-chat/domain"

type Messenger interface {
	ConnectClient(client *domain.Client)
	messageWriter(client *domain.Client, done chan<- struct{})
	messageReader(client *domain.Client, done <-chan struct{})
	updateActiveClients(client *domain.Client)
	disconnectClient(client *domain.Client, done chan<- struct{})
}
