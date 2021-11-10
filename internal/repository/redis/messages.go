package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepository(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client: client}
}

func (m MessagesRepository) PublishMessage(roomId string, message []byte) error {
	return m.client.Publish(context.Background(), roomId, message).Err()
}

func (m MessagesRepository) SubscribeToMessages(roomId string) *redis.PubSub {
	return m.client.Subscribe(context.Background(), roomId)
}
