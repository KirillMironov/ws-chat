package ports

import (
	"github.com/go-redis/redis/v8"
)

type ClientsRepository interface {
	PublishActiveClients(roomId string, message []byte) error
	SubscribeToActiveClients(roomId string) *redis.PubSub
	AddActiveClient(roomId, username string) error
	RemoveActiveClient(roomId, username string) error
	GetActiveClients(roomId string) ([]string, error)
}

type MessagesRepository interface {
	PublishMessage(roomId string, message []byte) error
	SubscribeToMessages(roomId string) *redis.PubSub
}
