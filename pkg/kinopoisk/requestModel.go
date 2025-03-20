package kinopoisk

type Response struct {
	Total      int    `json:"total"`
	TotalPages int    `json:"totalPages"`
	Items      []Item `json:"items"`
}

type Item struct {
	KinopoiskID int       `json:"kinopoiskId"`
	Countries   []Country `json:"countries"`
	Genres      []Genre   `json:"genres"`
}

type Country struct {
	Country string `json:"country"`
}

// Структура для представления жанра
type Genre struct {
	Genre string `json:"genre"`
}

type TrallerItem struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Site string `json:"site"`
}

// Структура триллера
type ResponseTraller struct {
	Total      int           `json:"total"`
	TotalPages int           `json:"totalPages"`
	Items      []TrallerItem `json:"items"`
}
