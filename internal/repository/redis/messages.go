package redis

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepository(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client: client}
}

func (m MessagesRepository) Publish(client *domain.Client, msg string) error {
	var message = domain.Message{Event: domain.ChatMessage}
	message.Payload.Username = client.Username
	message.Payload.Text = msg

	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return m.client.Publish(context.Background(), client.RoomId, encoded).Err()
}

func (m MessagesRepository) Subscribe(roomId string) *redis.PubSub {
	return m.client.Subscribe(context.Background(), roomId)
}
