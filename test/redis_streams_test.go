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
	roomId           = "main"
	usernameKey      = "username"
	numberOfMessages = 3
)

var counter int

func TestRedisStreams(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := subscribeToStream(&wg)
		if err != nil {
			t.Error(err)
		}
	}()

	wg.Wait()
	err := publishToStream()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, counter, numberOfMessages)
}

func publishToStream() error {
	for i := 0; i < numberOfMessages; i++ {
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

func subscribeToStream(wg *sync.WaitGroup) error {
	client.XGroupCreate(context.Background(), roomId, roomId, "0")
	wg.Done()

	for {
		result, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Streams: []string{roomId, ">"},
			Group:   roomId,
			Block:   time.Second * 10,
		}).Result()
		if err != nil {
			return err
		}

		counter++
		if result[0].Messages[0].Values[usernameKey] == strconv.Itoa(numberOfMessages-1) {
			return nil
		}
	}
}
