package main

import (
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/internal/delivery/v1"
	redisRepo "github.com/KirillMironov/ws-chat/internal/repository/redis"
	"github.com/KirillMironov/ws-chat/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func main() {
	// Logger
	logger := logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)
	logger.Level = logrus.DebugLevel

	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer client.Close()

	// App
	messagesRepo := redisRepo.NewMessagesRepo(client)
	messenger := service.NewWebSocketMessenger(messagesRepo, logger)
	handler := v1.NewHandler(messenger, logger)

	logger.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
