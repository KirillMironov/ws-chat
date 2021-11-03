package service

import "github.com/KirillMironov/ws-chat/domain"

type Messenger interface {
	ConnectClient(client *domain.Client)
	messageReader(client *domain.Client)
	messageWriter(client *domain.Client)
}
