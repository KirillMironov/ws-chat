package repository

import "github.com/go-redis/redis/v8"

type Messages interface {
	Publish(roomId string, message []byte) error
	Subscribe(roomId string) *redis.PubSub
}
