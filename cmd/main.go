package main

import (
	"context"
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/internal/delivery"
	repository "github.com/KirillMironov/ws-chat/internal/repository/redis"
	"github.com/KirillMironov/ws-chat/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

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
	messagesRepository := repository.NewMessagesRepository(client)
	clientsRepository := repository.NewClientsRepository(client)
	messagesService := service.NewMessagesService(messagesRepository, logger)
	clientsService := service.NewClientsService(clientsRepository, messagesService, logger)
	delivery.NewHandler(clientsService, logger).InitRoutes()

	logger.Infof("started on port %s", cfg.Port)
	logger.Fatal(http.ListenAndServe(":"+cfg.Port, http.DefaultServeMux))
}
