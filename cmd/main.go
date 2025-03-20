package main

import (
	"TelegramBot/config"
	"TelegramBot/internal/handler"
	"TelegramBot/internal/repository"
	"TelegramBot/internal/service"
	"TelegramBot/pkg/kinopoisk"
	"database/sql"
	"log"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.GetConfig()

	db, err := sql.Open("sqlite", "FilmsBook")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	bot, err := tgBotApi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	api := kinopoisk.NewKinopoisk(cfg.FilmsToken)
	repo := repository.NewRepository(db)
	service := service.NewFilms(repo, *api)
	handler := handler.NewHandler(service, *bot, *api)

	handler.InitHandler(1)
}
