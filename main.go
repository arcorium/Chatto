package main

import (
	"log"
	"os"

	"chatto/internal"
	"chatto/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := config.LoadConfig(workDir)
	if err != nil {
		log.Fatalln(err)
	}
	app := internal.NewApp(&cfg)

	app.Start()
}
