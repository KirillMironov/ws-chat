package redis

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/go-redis/redis/v8"
)

const keyPostfix = "clients"

type ClientsRepository struct {
	client *redis.Client
}

func NewClientsRepository(client *redis.Client) *ClientsRepository {
	return &ClientsRepository{client: client}
}

func (c ClientsRepository) Add(client *domain.Client) error {
	err := c.client.SAdd(context.Background(), client.RoomId+keyPostfix, client.Username).Err()
	if err != nil {
		return err
	}

	return c.publish(client.RoomId)
}

func (c ClientsRepository) Remove(client *domain.Client) error {
	err := c.client.SRem(context.Background(), client.RoomId+keyPostfix, client.Username).Err()
	if err != nil {
		return err
	}

	return c.publish(client.RoomId)
}

func (c ClientsRepository) publish(roomId string) error {
	clients, err := c.client.SMembers(context.Background(), roomId+keyPostfix).Result()
	if err != nil {
		return err
	}

	var message = domain.Message{Event: domain.ConnectedClients}
	message.Payload.Text = clients

	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.client.Publish(context.Background(), roomId, encoded).Err()
}
