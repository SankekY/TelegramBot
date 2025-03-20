package handler

import (
	"TelegramBot/internal/models"
	"TelegramBot/pkg/kinopoisk"
	"log"
	"time"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Films interface {
	SaveFilmsAndTrailers([]models.Film) error
	GetFilm() (string, []byte, int, error)
}

type Handler struct {
	filmService Films
	bot         tgBotApi.BotAPI
	api         kinopoisk.KinopoiskAPI
}

func NewHandler(films Films, bot tgBotApi.BotAPI, api kinopoisk.KinopoiskAPI) *Handler {
	return &Handler{filmService: films, bot: bot, api: api}
}

func (h *Handler) InitHandler(t int64) {
	go func() {
		h.postFilmToChanel()
		time.Sleep(time.Hour * time.Duration(t))
	}()
}

func (h *Handler) postFilmToChanel() {
	caption, file, NotPosted, err := h.filmService.GetFilm()
	if err != nil {
		log.Println("Error :", err)
		return
	}
	if NotPosted <= 10 {
		go h.saveFilmsAndTrailers()
	}
	newPhoto := tgBotApi.PhotoConfig{
		Thumb:     tgBotApi.FileBytes{Bytes: file},
		Caption:   caption,
		ParseMode: "html",
	}
	newPhoto.ChannelUsername = "Channel Name"
	h.bot.Send(newPhoto)
}

func (h *Handler) saveFilmsAndTrailers() {
	films := h.api.GetStackFilms()
	err := h.filmService.SaveFilmsAndTrailers(films)
	if err != nil {
		log.Println("Error : ", err)
	}
}
