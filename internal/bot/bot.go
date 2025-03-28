package bot

import (
	"TelegramBot/internal/models"
	"TelegramBot/internal/service"
	"TelegramBot/pkg/kinopoisk"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	bot         *tgBotApi.BotAPI
	filmService service.Films
	api         *kinopoisk.KinopoiskAPI
}

var msgStart = `✨ Команды бота ✨  

🎬 Добавить фильм  
/add film_id – добавь фильм в избранное  

🗑 Удалить фильм  
/del film_id – удали фильм из избранного  

📋 Мои фильмы  
/my – посмотри свой список избранного  

📢 Опубликовать в канал  
/post film_id – отправь фильм в канал  

🙏 NAMASTE  
*Меньше чем три...* 💫`

func New(bot *tgBotApi.BotAPI, films service.Films, api *kinopoisk.KinopoiskAPI) *tgBot {
	return &tgBot{bot: bot, filmService: films, api: api}
}

func (t *tgBot) BotInit() {
	update := tgBotApi.NewUpdate(1)
	updates := t.bot.GetUpdatesChan(update)
	for u := range updates {
		msg := tgBotApi.NewMessage(u.Message.Chat.ID, "")
		msg.ParseMode = "html"
		if u.Message != nil {
			spltMsg := strings.Split(u.Message.Text, " ")
			if len(spltMsg) > 1 {
				switch spltMsg[0] {
				case "/add":
					filmId, err := strconv.Atoi(spltMsg[1])
					if err != nil {
						msg.Text = "Пришлите валиднное ID !"
						t.bot.Send(msg)
						continue
					}
					item, err := t.api.GetFullInfoByID(filmId)
					if err != nil {
						msg.Text = "По данному ID фильма не найдено!\n Прешлите ID с кинопоиска !"
						t.bot.Send(msg)
						continue
					}
					if err := t.filmService.SaveFilmToUserPool(models.FilmsPool{
						UserID:        u.Message.From.ID,
						KinopoiskID:   item.KinopoiskID,
						KinopoiskHDID: item.KinopoiskHDID,
						Title:         item.TitleRu,
						Year:          item.Year,
						Description:   item.Description,
						PosterUrl:     item.PosterUrl,
					}); err != nil {
						msg.Text = "Bad Request: :( Фильм не сохранён попробуйте позже :( "
						log.Println(err)
						t.bot.Send(msg)
						continue
					}
					msg.Text = fmt.Sprintf("Фильм: %s %d\n Сохранён в избранное!", item.TitleRu, item.Year)
					t.bot.Send(msg)
				case "/del":
					filmId, err := strconv.Atoi(spltMsg[1])
					if err != nil {
						msg.Text = "Пришлите валиднное ID !"
						t.bot.Send(msg)
						continue
					}
					if err := t.filmService.DeleteFilmUserPool(filmId, u.Message.From.ID); err != nil {
						msg.Text = "Bad Request: %( Не верный ID или ошибка на сервере %("
						log.Println(err)
						t.bot.Send(msg)
						continue
					}
					msg.Text = "Ок: Фильм удалён !"
					t.bot.Send(msg)
				case "/post":
					filmId, err := strconv.Atoi(spltMsg[1])
					if err != nil {
						msg.Text = "Пришлите валиднное ID !"
						t.bot.Send(msg)
						continue
					}
					caption, file, err := t.filmService.UserPost(filmId)
					if err != nil {
						msg.Text = "Bad Request"
						t.bot.Send(msg)
						continue
					}
					newPhoto := tgBotApi.NewPhotoToChannel("@cahnel_name", tgBotApi.FileBytes{
						Name:  "name",
						Bytes: file,
					})
					newPhoto.Caption = caption
					newPhoto.ParseMode = "html"
					t.bot.Send(newPhoto)
				}
			}
			switch u.Message.Text {
			case "/start":
				t.filmService.SaveUser(models.User{
					UserID:   u.Message.From.ID,
					ChatID:   u.Message.Chat.ID,
					UserName: u.Message.From.UserName,
				})
				msg.Text = msgStart
				t.bot.Send(msg)
			case "/my":
				text, err := t.filmService.GetUserPool(
					u.Message.From.ID,
					"false",
				)
				if err != nil {
					log.Println("bot.BotInit.mylist: ", err)
					msg.Text = "Тех неполадки или у вас нету фильмов в пуле :("
					t.bot.Send(msg)
					continue
				}
				msg.Text = text

				t.bot.Send(msg)
			}
		}
	}
}
