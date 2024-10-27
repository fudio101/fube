package main

import (
	"github.com/fudio101/fube/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg, errGetConfig := config.LoadConfig()
	if errGetConfig != nil {
		log.Fatalf("Failed to load config: %v", errGetConfig)

	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	errRunEngine := r.Run(cfg.ServerAddress)
	if errRunEngine != nil {
		log.Fatalf("Failed to run server: %v", errRunEngine)
	}
}
