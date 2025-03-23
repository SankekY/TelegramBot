package service

import (
	"TelegramBot/internal/models"
	"TelegramBot/pkg/kinopoisk"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type FilmsRepository interface {
	GetFilm() (models.Film, []models.Trailer, error)
	SaveFilmsAndTrailers(films []models.Film)
	SaveUser(user models.User)
	SaveFilmToUserPool(film models.FilmsPool) error
	GetUserPool(userId int64, watch string) ([]models.FilmsPool, error)
	DeleteFilmUserPool(filmId int, userId int64) error
}

type Films struct {
	repo      FilmsRepository
	kinopoisk kinopoisk.KinopoiskAPI
}

func NewFilms(repo FilmsRepository, kino kinopoisk.KinopoiskAPI) *Films {
	return &Films{repo, kino}
}

func (f *Films) GetFilm() (string, []byte, int, error) {
	film, trailers, err := f.repo.GetFilm()
	if err != nil {
		log.Println(err)
		return "", nil, 0, err
	}
	poserByte, err := getPosterByte(film.PosterUrl)
	if err != nil {
		return "", nil, 0, err
	}
	posText := generateText(film, trailers)

	return posText, poserByte, film.NonPosted, nil
}

func (f *Films) SaveFilmsAndTrailers(films []models.Film) {
	f.repo.SaveFilmsAndTrailers(films)
}

func (f *Films) SaveUser(user models.User) {
	f.repo.SaveUser(user)
}

func (f *Films) SaveFilmToUserPool(film models.FilmsPool) error {
	return f.repo.SaveFilmToUserPool(film)
}

func (f *Films) GetUserPool(userId int64, watch string) (string, error) {
	films, err := f.repo.GetUserPool(userId, watch)
	if err != nil {
		return "", err
	}
	text := generateTextForUser(films)
	return text, nil
}

func (f *Films) UserPost(filmId int) (string, []byte, error) {
	film, err := f.kinopoisk.GetFullInfoByID(filmId)
	if err != nil {
		return "", nil, err
	}
	film.Trailers = f.kinopoisk.GetReqTrailers(film.KinopoiskID)
	text := generateText(film, film.Trailers)
	file, err := getPosterByte(film.PosterUrl)
	if err != nil {
		return text, nil, nil
	}
	return text, file, nil
}

func (f *Films) DeleteFilmUserPool(filmId int, userId int64) error {
	return f.repo.DeleteFilmUserPool(filmId, userId)
}

func generateTextForUser(films []models.FilmsPool) string {
	text := "Ваша библеотека фильмов! \n\n"
	for _, film := range films {
		freeUrl := fmt.Sprintf("https://r.kpfr.site/film/%d/", film.KinopoiskID)
		text += fmt.Sprintf("<b>%s | %d</b>\n", film.Title, film.Year)
		text += fmt.Sprintf("<b>Watch:</b> <a href='%s'>FeeWatch</a>\n<b>FilmID-->:</b> <i>%d</i>\n\n ", freeUrl, film.KinopoiskID)
	}

	return text
}

func getPosterByte(url string) ([]byte, error) {
	resp, _ := http.Get(url)
	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func generateText(film models.Film, tralelers []models.Trailer) string {
	text := fmt.Sprintf("%s | %s \n Год: %d | Рейтинг: %v #%s\n%s\n",
		film.TitleRu, film.TitleOrig, film.Year, film.Rating, film.Genre,
		film.Description,
	)
	PoiskUrl := fmt.Sprintf("https://hd.kinopoisk.ru/film/%s", film.KinopoiskHDID)
	freeUrl := fmt.Sprintf("https://r.kpfr.site/film/%d/", film.KinopoiskID)

	youtube := 0
	yandex := 0
	for _, v := range tralelers {
		if v.Site == "YOUTUBE" {
			if youtube == 0 {
				text += fmt.Sprintf("YouTube: <a href='%s'>%s</a> \n", v.URL, v.Title)
				youtube++
			}

		}
		if v.Site == "YANDEX_DISK" {
			if yandex == 0 {
				text += fmt.Sprintf("Yandex: href='%s'>%s</a> \n", v.URL, v.Title)
				yandex++
			}
		}
	}
	text += fmt.Sprintf("Watch: <a href='%s'>KinopoiskHD</a>\n", PoiskUrl)
	text += fmt.Sprintf("WatchFree: <a href='%s'>Thanks</a>\n", freeUrl)
	text += fmt.Sprintf("filmID-->: %d", film.KinopoiskID)
	return text
}
