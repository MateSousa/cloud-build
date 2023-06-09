package main

import (
	"log"

	"github.com/MateSousa/cloud-build/cmd/api"
	"github.com/MateSousa/cloud-build/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
	initializers.ConnectRedis(&config)

	api.Init()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	log.Default().Println(config)

}
