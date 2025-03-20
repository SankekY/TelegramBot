package models

type Film struct {
	KinopoiskID   int     `json:"kinopoiskId" db:"KinopoiskID"`
	KinopoiskHDID string  `json:"kinopoiskHDId" db:"KinopoiskHDID"`
	TitleRu       string  `json:"nameRu" db:"TitleRu"`
	TitleOrig     string  `json:"nameOriginal" db:"TitleOrig"`
	PosterUrl     string  `json:"posterUrl" db:"PosterUrl"`
	Description   string  `json:"short_description" db:"Description"`
	Country       string  `json:"country" db:"Country"`
	Year          int     `json:"year" db:"Year"`
	Rating        float64 `json:"rating" db:"Rating"`
	Genre         string  `json:"genre" db:"Genre"`
	Posted        string  `json:"posted" db:"Posted"`
	NonPosted     int     `db:"filmPosted"`
	Trailers      []Trailer
}

type Trailer struct {
	KinopoiskID int    `json:"kinopoisk_id" db:"kinopoiskID"`
	URL         string `json:"url" db:"url"`
	Title       string `json:"name" db:"title"`
	Site        string `json:"site" db:"site"`
}

type User struct {
	UserID   int    `json:"user_id" db:"user_id"`
	ChatID   int    `json:"chat_id" db:"chat_id"`
	UserName string `json:"user_name" db:"user_name"`
}

type FilmsPool struct {
	UserID      int    `json:"user_id" db:"user_id"`
	KinopoiskID int    `json:"kinopoisk_id" db:"kinopoisk_id"`
	Year        int    `json:"year" db:"year"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	PosterUrl   string `json:"poster_url" db:"poster_url"`
	Watch       int    `json:"watch" db:"watch"`
}
