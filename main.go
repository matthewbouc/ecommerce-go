package main

import (
	"ecommerce/config"
	"ecommerce/internal/api"
	"log"
)

func main() {

	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalf("config file did not load properly - %v\n", err)
	}

	api.StartServer(cfg)
}
