package test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

const (
	roomId     = "main"
	username   = "Lisa"
	keyPostfix = ":clients"
)

func TestAddActiveUser(t *testing.T) {
	t.Skip()

	err := client.SAdd(context.Background(), roomId+keyPostfix, username).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveActiveUser(t *testing.T) {
	t.Skip()

	err := client.SRem(context.Background(), roomId+keyPostfix, username).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetActiveUsers(t *testing.T) {
	t.Skip()

	result, err := client.SMembers(context.Background(), roomId+keyPostfix).Result()
	if err != nil {
		t.Error(err)
	}
	log.Println(result)
}

func TestPublishActiveUsers(t *testing.T) {
	t.Skip()

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
	t.Skip()

	subscription := client.Subscribe(context.Background(), roomId)
	defer subscription.Close()

	for msg := range subscription.Channel() {
		log.Println(msg)
	}
}

func getActiveUsers(key string) ([]string, error) {
	result, err := client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
