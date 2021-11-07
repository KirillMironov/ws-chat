package test

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	roomId          = "main"
	usernameKey     = "username"
	numberOfMessage = 1000
)

var results []interface{}

func TestRedisStreams(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(t *testing.T) {
		defer wg.Done()
		err := readFromStream()
		if err != nil {
			t.Error(err)
		}
	}(t)

	go func(t *testing.T) {
		defer wg.Done()
		err := publishToStream()
		if err != nil {
			t.Error(err)
		}
	}(t)

	wg.Wait()

	assert.Equal(t, 1000, len(results))
}

func publishToStream() error {
	for i := 0; i < numberOfMessage; i++ {
		err := client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: roomId,
			Values: map[string]interface{}{usernameKey: strconv.Itoa(i)},
		}).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func readFromStream() error {
	client.XGroupCreate(context.Background(), roomId, roomId, "0")

	for {
		result, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Streams: []string{roomId, ">"},
			Group:   roomId,
			Block:   time.Second * 10,
		}).Result()
		if err != nil {
			return err
		}

		if result != nil {
			results = append(results, result[0].Messages[0].Values[usernameKey])
		}

		if result[0].Messages[0].Values[usernameKey] == strconv.Itoa(numberOfMessage-1) {
			return nil
		}
	}
}
