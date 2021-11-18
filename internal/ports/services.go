package ports

import (
	"github.com/KirillMironov/ws-chat/domain"
)

type ClientsService interface {
	Connect(client *domain.Client)
	Disconnect(client *domain.Client, done chan<- struct{})
}

type MessagesService interface {
	Reader(client *domain.Client)
	Writer(client *domain.Client, done chan struct{})
}
