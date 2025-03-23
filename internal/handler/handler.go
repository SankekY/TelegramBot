package handler

import (
	"TelegramBot/internal/models"
	"TelegramBot/pkg/kinopoisk"
	"log"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Films interface {
	SaveFilmsAndTrailers([]models.Film)
	GetFilm() (string, []byte, int, error)
	SaveUser(user models.User)
	SaveFilmToUserPool(film models.FilmsPool) error
	GetUserPool(userId int64, watch string) (string, error)
}

type Handler struct {
	filmService Films
	bot         tgBotApi.BotAPI
	api         kinopoisk.KinopoiskAPI
}

func NewHandler(films Films, bot tgBotApi.BotAPI, api kinopoisk.KinopoiskAPI) *Handler {
	return &Handler{filmService: films, bot: bot, api: api}
}

func (h *Handler) InitHandler() {
	go func() {

		// h.postFilmToChanel()
	}()
}

func (h *Handler) postFilmToChanel() {
	caption, file, sum, err := h.filmService.GetFilm()
	if sum <= 10 {
		// h.saveFilmsAndTrailers()
	}
	if err != nil {
		log.Println(err)
		return
	}

	newPhoto := tgBotApi.NewPhotoToChannel("@cahnel_name", tgBotApi.FileBytes{
		Name:  "name",
		Bytes: file,
	})
	newPhoto.Caption = caption
	newPhoto.ParseMode = "html"
	h.bot.Send(newPhoto)
}

func (h *Handler) saveFilmsAndTrailers() {
	log.Println("Go save films")
	films := h.api.GetStackFilms()
	h.filmService.SaveFilmsAndTrailers(films)

}
