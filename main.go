package main

import (
	"LoRaPanel/backend"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("[STARTUP] Starting LoRaPanel server instance...")
	godotenv.Load()

	server := backend.InitBackend(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DATABASE"))
	err := server.StartBackend()
	log.Fatalf("Server Fatal Error %v", err)
}