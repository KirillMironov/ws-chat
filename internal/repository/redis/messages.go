package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const usersCounterPostfix = ":activeUsers"

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

func (m MessagesRepository) AddActiveUser(roomId, username string) error {
	return m.client.SAdd(context.Background(), roomId+usersCounterPostfix, username).Err()
}

func (m MessagesRepository) RemoveActiveUser(roomId, username string) error {
	return m.client.SRem(context.Background(), roomId+usersCounterPostfix, username).Err()
}

func (m MessagesRepository) GetActiveUsers(roomId string) ([]string, error) {
	return m.client.SMembers(context.Background(), roomId+usersCounterPostfix).Result()
}

func (m MessagesRepository) PublishActiveUsers(roomId string, message []byte) error {
	return m.client.Publish(context.Background(), roomId+usersCounterPostfix, message).Err()
}

func (m MessagesRepository) SubscribeToActiveUsers(roomId string) *redis.PubSub {
	return m.client.Subscribe(context.Background(), roomId+usersCounterPostfix)
}
