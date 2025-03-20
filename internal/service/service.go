package service

import (
	"TelegramBot/internal/models"
	"TelegramBot/pkg/kinopoisk"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FilmsRepository interface {
	GetFilm() (models.Film, []models.Trailer, error)
	SaveFilmsAndTrailers(films []models.Film) error
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
		return "", nil, 0, err
	}
	poserByte, err := getPosterByte(film.PosterUrl)
	if err != nil {
		return "", nil, 0, err
	}
	posText := generateText(film, trailers)
	if film.NonPosted <= 6 {

	}
	return posText, poserByte, film.NonPosted, nil
}

func (f *Films) SaveFilmsAndTrailers(films []models.Film) error {
	return f.repo.SaveFilmsAndTrailers(films)
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
	text += fmt.Sprintf("Watch: <a herf='%s'>KinopoiskHD</a>\n", PoiskUrl)
	text += fmt.Sprintf("WatchFree: <a herf='%s'>Thanks</a>\n", freeUrl)
	youtube := 0
	yandex := 0
	for _, v := range tralelers {
		if v.Site == "YOUTUBE" && youtube < 2 {
			text += fmt.Sprintf("YouTube: <a herf='%s'>%s</a> \n", v.URL, v.Title)
		}
		if v.Site == "YANDEX_DISK" && yandex < 2 {
			text += fmt.Sprintf("Yandex: <a herf='%s'>%s</a> \n", v.URL, v.Title)
		}
	}
	return text
}
