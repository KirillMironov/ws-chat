package repository

import "github.com/go-redis/redis/v8"

type Messages interface {
	SendMessage(roomId string, message []byte) error
	GetMessages(roomId string) <-chan *redis.Message
}
