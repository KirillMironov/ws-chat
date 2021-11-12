package test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
)

var client *redis.Client

func TestMain(m *testing.M) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
