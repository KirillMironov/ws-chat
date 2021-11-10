package ports

import (
	"github.com/KirillMironov/ws-chat/domain"
	"sync"
)

type ClientsService interface {
	ConnectClient(client *domain.Client)
	DisconnectClient(client *domain.Client, done chan<- struct{})
	UpdateActiveClients(client *domain.Client)
}

type MessagesService interface {
	MessageWriter(client *domain.Client)
	MessageReader(client *domain.Client, done <-chan struct{}, wg *sync.WaitGroup)
}
