package test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

const (
	keyPostfix = ":activeUsers"
	roomId     = "main"
	username   = "Lisa"
)

func TestAddActiveUser(t *testing.T) {
	err := client.SAdd(context.Background(), roomId+keyPostfix, username).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveActiveUser(t *testing.T) {
	err := client.SRem(context.Background(), roomId+keyPostfix, username).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetActiveUsers(t *testing.T) {
	result, err := client.SMembers(context.Background(), roomId+keyPostfix).Result()
	if err != nil {
		t.Error(err)
	}
	log.Println(result)
}

func TestPublishActiveUsers(t *testing.T) {
	users, err := getActiveUsers(roomId + keyPostfix)
	if err != nil {
		t.Error(err)
	}

	js, err := json.Marshal(users)
	if err != nil {
		t.Error(err)
	}

	err = client.Publish(context.Background(), roomId+keyPostfix, js).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestSubscribeToActiveUsers(t *testing.T) {
	subscription := client.Subscribe(context.Background(), roomId+keyPostfix)
	defer subscription.Unsubscribe(context.Background(), roomId+keyPostfix)

	for {
		select {
		case msg := <-subscription.Channel():
			log.Println(msg)
		}
	}
}

func getActiveUsers(key string) ([]string, error) {
	result, err := client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
