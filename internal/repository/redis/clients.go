package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type ClientsRepository struct {
	client *redis.Client
}

func NewClientsRepository(client *redis.Client) *ClientsRepository {
	return &ClientsRepository{client: client}
}

const connectedClientsPostfix = ":connectedClients"

func (c ClientsRepository) Publish(roomId string, message []byte) error {
	return c.client.Publish(context.Background(), roomId+connectedClientsPostfix, message).Err()
}

func (c ClientsRepository) Subscribe(roomId string) *redis.PubSub {
	return c.client.Subscribe(context.Background(), roomId+connectedClientsPostfix)
}

func (c ClientsRepository) Add(roomId, username string) error {
	return c.client.SAdd(context.Background(), roomId+connectedClientsPostfix, username).Err()
}

func (c ClientsRepository) Remove(roomId, username string) error {
	return c.client.SRem(context.Background(), roomId+connectedClientsPostfix, username).Err()
}

func (c ClientsRepository) GetConnected(roomId string) ([]string, error) {
	return c.client.SMembers(context.Background(), roomId+connectedClientsPostfix).Result()
}
