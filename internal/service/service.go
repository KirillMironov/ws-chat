package service

import (
	"github.com/KirillMironov/ws-chat/domain"
	"sync"
)

type Clients interface {
	ConnectClient(client *domain.Client)
	disconnectClient(client *domain.Client, done chan<- struct{})
	updateActiveClients(client *domain.Client)
}

type Messages interface {
	MessageWriter(client *domain.Client)
	MessageReader(client *domain.Client, done <-chan struct{}, wg *sync.WaitGroup)
}
