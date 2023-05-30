package main

import (
	"fmt"
	"log"

	"github.com/MateSousa/cloud-build/initializers"
	"github.com/MateSousa/cloud-build/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Token{})
	
	fmt.Println("? Migration complete")
}

