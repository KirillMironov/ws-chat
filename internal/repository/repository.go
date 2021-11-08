package repository

import (
	"github.com/go-redis/redis/v8"
)

type Messages interface {
	PublishMessage(roomId string, message []byte) error
	SubscribeToMessages(roomId string) *redis.PubSub
	AddActiveUser(roomId, username string) error
	RemoveActiveUser(roomId, username string) error
	GetActiveUsers(roomId string) ([]string, error)
	PublishActiveUsers(roomId string, message []byte) error
	SubscribeToActiveUsers(roomId string) *redis.PubSub
}
