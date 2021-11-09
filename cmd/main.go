package main

import (
	"context"
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/internal/delivery/v1"
	_repo "github.com/KirillMironov/ws-chat/internal/repository/redis"
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

	err = client.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal(err)
	}

	// App
	clientsRepo := _repo.NewClientsRepository(client)
	messagesRepo := _repo.NewMessagesRepository(client)
	messagesService := service.NewMessagesService(clientsRepo, messagesRepo, logger)
	clientsService := service.NewClientsService(clientsRepo, messagesService, logger)
	handler := v1.NewHandler(clientsService, logger)

	logger.Infof("started on port %s", cfg.Port)
	logger.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
