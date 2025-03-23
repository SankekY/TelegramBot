package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	FilmsToken  string
	ChannelName string
}

func GetConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Warrning by load config: ", err)
	}
	log.Println("Config Loaded")
	return &Config{
		BotToken:    os.Getenv("telegram_api_token"),
		FilmsToken:  os.Getenv("kinopoisk_api_token"),
		ChannelName: os.Getenv("channel_name"),
	}
}
