package ports

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

type ClientsRepository interface {
	Add(client *domain.Client) error
	Remove(client *domain.Client) error
	Publish(roomId string) error
}

type MessagesRepository interface {
	Publish(client *domain.Client, message []byte) error
	Subscribe(roomId string) *redis.PubSub
}
