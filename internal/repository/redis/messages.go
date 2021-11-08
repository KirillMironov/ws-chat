package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const activeClientsPostfix = ":activeClients"

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepo(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client: client}
}

func (m MessagesRepository) PublishMessage(roomId string, message []byte) error {
	return m.client.Publish(context.Background(), roomId, message).Err()
}

func (m MessagesRepository) SubscribeToMessages(roomId string) *redis.PubSub {
	return m.client.Subscribe(context.Background(), roomId)
}

func (m MessagesRepository) PublishActiveClients(roomId string, message []byte) error {
	return m.client.Publish(context.Background(), roomId+activeClientsPostfix, message).Err()
}

func (m MessagesRepository) SubscribeToActiveClients(roomId string) *redis.PubSub {
	return m.client.Subscribe(context.Background(), roomId+activeClientsPostfix)
}

func (m MessagesRepository) AddActiveClient(roomId, username string) error {
	return m.client.SAdd(context.Background(), roomId+activeClientsPostfix, username).Err()
}

func (m MessagesRepository) RemoveActiveClient(roomId, username string) error {
	return m.client.SRem(context.Background(), roomId+activeClientsPostfix, username).Err()
}

func (m MessagesRepository) GetActiveClients(roomId string) ([]string, error) {
	return m.client.SMembers(context.Background(), roomId+activeClientsPostfix).Result()
}
