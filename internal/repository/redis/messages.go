package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepo(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client: client}
}

func (m MessagesRepository) SendMessage(roomId string, message []byte) error {
	return m.client.Publish(context.Background(), roomId, message).Err()
}

func (m MessagesRepository) GetMessages(roomId string) <-chan *redis.Message {
	sub := m.client.Subscribe(context.Background(), roomId)
	return sub.Channel()
}
