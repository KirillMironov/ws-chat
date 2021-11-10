package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const activeClientsPostfix = ":activeClients"

type ClientsRepository struct {
	client *redis.Client
}

func NewClientsRepository(client *redis.Client) *ClientsRepository {
	return &ClientsRepository{client: client}
}

func (c ClientsRepository) PublishActiveClients(roomId string, message []byte) error {
	return c.client.Publish(context.Background(), roomId+activeClientsPostfix, message).Err()
}

func (c ClientsRepository) SubscribeToActiveClients(roomId string) *redis.PubSub {
	return c.client.Subscribe(context.Background(), roomId+activeClientsPostfix)
}

func (c ClientsRepository) AddActiveClient(roomId, username string) error {
	return c.client.SAdd(context.Background(), roomId+activeClientsPostfix, username).Err()
}

func (c ClientsRepository) RemoveActiveClient(roomId, username string) error {
	return c.client.SRem(context.Background(), roomId+activeClientsPostfix, username).Err()
}

func (c ClientsRepository) GetActiveClients(roomId string) ([]string, error) {
	return c.client.SMembers(context.Background(), roomId+activeClientsPostfix).Result()
}
