package main

import (
	"log"
	"os"

	"server_client_chat/internal"
	"server_client_chat/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	cfg, err := config.LoadConfig(workDir)
	if err != nil {
		log.Fatalln(err)
	}
	server := internal.NewApp(&cfg)

	server.Start()
}
