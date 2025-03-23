package models

type Film struct {
	KinopoiskID   int     `json:"kinopoiskId" db:"KinopoiskID"`
	KinopoiskHDID string  `json:"kinopoiskHDId" db:"KinopoiskHDID"`
	TitleRu       string  `json:"nameRu" db:"TitleRu"`
	TitleOrig     string  `json:"nameOriginal" db:"TitleOrig"`
	PosterUrl     string  `json:"posterUrl" db:"PosterUrl"`
	Description   string  `json:"description" db:"Description"`
	Country       string  `json:"country" db:"Country"`
	Year          int     `json:"year" db:"Year"`
	Rating        float64 `json:"ratingKinopoisk" db:"Rating"`
	Genre         string  `json:"genre" db:"Genre"`
	Posted        string  `json:"posted" db:"Posted"`
	NonPosted     int     `db:"filmPosted"`
	Favor         int     `db:"Favor"`
	Trailers      []Trailer
}

type Trailer struct {
	KinopoiskID int    `json:"kinopoisk_id" db:"kinopoiskID"`
	URL         string `json:"url" db:"url"`
	Title       string `json:"name" db:"title"`
	Site        string `json:"site" db:"site"`
}

type User struct {
	UserID   int64  `json:"user_id" db:"UserId"`
	ChatID   int64  `json:"chat_id" db:"ChatId"`
	UserName string `json:"user_name" db:"UserName"`
}

type FilmsPool struct {
	UserID        int64  `json:"user_id" db:"UserID"`
	KinopoiskID   int    `json:"kinopoisk_id" db:"KinpoiskID"`
	KinopoiskHDID string `db:"KinopoiskHDID"`
	Year          int    `json:"year" db:"Year"`
	Title         string `json:"title" db:"Title"`
	Description   string `json:"description" db:"Description"`
	PosterUrl     string `json:"poster_url" db:"PosterUrl"`
	Watch         string `json:"watch" db:"Watch"`
}
