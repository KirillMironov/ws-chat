package main

import (
	"github.com/KirillMironov/ws-chat/config"
	"github.com/KirillMironov/ws-chat/internal/delivery"
	"log"
)

func main() {
	cfg := config.InitConfig()
	log.Fatal(delivery.InitRoutes().Run(":" + cfg.Port))
}
