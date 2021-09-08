package main

import (
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/delivery"
	"github.com/KirillMironov/ws-chat/internal/service"
	"log"
)

func main() {
	cfg := config.InitConfig()
	rooms := make(map[string]domain.Room)

	messenger := service.NewMessengerWS(rooms)
	handler := delivery.NewHandler(messenger)
	log.Fatal(handler.InitRoutes().Run(":" + cfg.Port))
}
