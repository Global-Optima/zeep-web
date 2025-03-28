package main

import (
	"fmt"
	"log"

	initializer "github.com/Global-Optima/zeep-web/backend/cmd/init"
)

func main() {
	router, cfg := initializer.InitializeApp()

	port := cfg.Server.Port
	log.Printf("Starting server on port %d...", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
