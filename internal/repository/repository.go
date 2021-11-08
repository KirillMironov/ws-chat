package repository

import (
	"github.com/go-redis/redis/v8"
)

type Messages interface {
	PublishMessage(roomId string, message []byte) error
	SubscribeToMessages(roomId string) *redis.PubSub
	PublishActiveClients(roomId string, message []byte) error
	SubscribeToActiveClients(roomId string) *redis.PubSub
	AddActiveClient(roomId, username string) error
	RemoveActiveClient(roomId, username string) error
	GetActiveClients(roomId string) ([]string, error)
}
