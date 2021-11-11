package ports

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

type ClientsRepository interface {
	Publish(client *domain.Client, message []byte) error
	Subscribe(client *domain.Client) *redis.PubSub
	Add(client *domain.Client) error
	Remove(client *domain.Client) error
	GetConnected(client *domain.Client) ([]string, error)
}

type MessagesRepository interface {
	Publish(client *domain.Client, message []byte) error
	Subscribe(client *domain.Client) *redis.PubSub
}
