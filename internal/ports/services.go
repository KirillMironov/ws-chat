package ports

import (
	"github.com/KirillMironov/ws-chat/domain"
	"sync"
)

type ClientsService interface {
	Connect(client *domain.Client)
	Disconnect(client *domain.Client, done chan<- struct{})
	UpdateConnected(client *domain.Client)
}

type MessagesService interface {
	Reader(client *domain.Client)
	Writer(client *domain.Client, done <-chan struct{}, wg *sync.WaitGroup)
}
