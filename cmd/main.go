package main

import (
	"TelegramBot/config"
	"TelegramBot/internal/bot"
	"TelegramBot/internal/handler"
	"TelegramBot/internal/repository"
	"TelegramBot/internal/service"
	"TelegramBot/pkg/kinopoisk"
	"database/sql"
	"log"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.GetConfig()
	db, err := sql.Open("sqlite3", "FilmsBook.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	tgBot, err := tgBotApi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	tgBot.Debug = false

	api := kinopoisk.NewKinopoisk(cfg.FilmsToken)
	repo := repository.NewRepository(db)
	service := service.NewFilms(repo, *api)
	handler := handler.NewHandler(service, *tgBot, *api)
	b := bot.New(tgBot, *service, api)
	handler.InitHandler()
	b.BotInit()

	ch := make(chan int)
	<-ch
}
