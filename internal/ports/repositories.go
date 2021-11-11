package ports

import (
	"github.com/go-redis/redis/v8"
)

type ClientsRepository interface {
	Publish(roomId string, message []byte) error
	Subscribe(roomId string) *redis.PubSub
	Add(roomId, username string) error
	Remove(roomId, username string) error
	GetConnected(roomId string) ([]string, error)
}

type MessagesRepository interface {
	Publish(roomId string, message []byte) error
	Subscribe(roomId string) *redis.PubSub
}
