package redis

import (
	"context"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

type ClientsRepository struct {
	client *redis.Client
}

func NewClientsRepository(client *redis.Client) *ClientsRepository {
	return &ClientsRepository{client: client}
}

const connectedClientsPostfix = ":connectedClients"

func (c ClientsRepository) Publish(client *domain.Client, message []byte) error {
	return c.client.Publish(context.Background(), client.RoomId+connectedClientsPostfix, message).Err()
}

func (c ClientsRepository) Subscribe(client *domain.Client) *redis.PubSub {
	return c.client.Subscribe(context.Background(), client.RoomId+connectedClientsPostfix)
}

func (c ClientsRepository) Add(client *domain.Client) error {
	return c.client.SAdd(context.Background(), client.RoomId+connectedClientsPostfix, client.Username).Err()
}

func (c ClientsRepository) Remove(client *domain.Client) error {
	return c.client.SRem(context.Background(), client.RoomId+connectedClientsPostfix, client.Username).Err()
}

func (c ClientsRepository) GetConnected(client *domain.Client) ([]string, error) {
	return c.client.SMembers(context.Background(), client.RoomId+connectedClientsPostfix).Result()
}
