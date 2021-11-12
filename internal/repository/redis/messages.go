package redis

import (
	"context"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepository(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client: client}
}

func (m MessagesRepository) Publish(client *domain.Client, message []byte) error {
	return m.client.Publish(context.Background(), client.RoomId, message).Err()
}

func (m MessagesRepository) Subscribe(client *domain.Client) *redis.PubSub {
	return m.client.Subscribe(context.Background(), client.RoomId)
}
