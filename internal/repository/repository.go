package repository

import "github.com/go-redis/redis/v8"

type Messages interface {
	SendMessage(roomId, message string) error
	GetMessages(roomId string) <-chan *redis.Message
}
