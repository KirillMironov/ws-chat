package service

import "github.com/KirillMironov/ws-chat/domain"

type Messenger interface {
	ConnectClient(client *domain.Client)
	messageWriter(client *domain.Client, closeSignal chan<- struct{})
	messageReader(client *domain.Client, closeSignal <-chan struct{})
}
