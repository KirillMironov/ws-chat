package main

import (
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/delivery"
	"github.com/KirillMironov/ws-chat/internal/service"
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

	// App
	rooms := make(map[string]domain.Room)

	messenger := service.NewWebSocketMessenger(rooms, logger)
	handler := delivery.NewHandler(messenger, logger)

	logger.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
