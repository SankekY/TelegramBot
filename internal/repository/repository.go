package repository

import (
	"TelegramBot/internal/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type FilmRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{db}
}

func (repo *FilmRepository) GetFilm() (models.Film, []models.Trailer, error) {
	const op = "repository.GetFilm"
	film, err := scannerFilm(repo.db.QueryRow("SELECT * FROM Films WHERE Posted = 'false'"))
	if err != nil {
		log.Printf("Erorr [%s], err: %s", op, err)
		return models.Film{}, nil, fmt.Errorf("failed to retrieve film: %w", err)
	}

	rows, err := repo.db.Query("SELECT * FROM Trailers WHERE KinopoiskID = ?", film.KinopoiskID)
	if err != nil {
		return film, nil, fmt.Errorf("failed to retrieve trailers: %w", err)
	}
	defer rows.Close()
	trailers, err := scannerTrailer(rows)
	if err != nil {
		return film, nil, nil
	}
	repo.db.Exec("update Films set Posted = 'true' WHERE KinopoiskID = ? ", film.KinopoiskID)

	repo.db.QueryRow("SELECT COUNT(*) AS filmPosted FROM Films WHERE Posted = 'false'").Scan(&film.NonPosted)

	return film, trailers, nil
}

func (repo *FilmRepository) SaveFilmsAndTrailers(films []models.Film) {
	for _, film := range films {
		_, err := repo.db.Exec(
			"INSERT INTO Films (KinopoiskID, KinopoiskHDID, TitleRu, TitleOrig, PosterUrl, Description, Country, Year, Rating, Genre) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			film.KinopoiskID, film.KinopoiskHDID, film.TitleRu, film.TitleOrig,
			film.PosterUrl, film.Description, film.Country,
			film.Year, film.Rating, film.Genre,
		)
		if err != nil {
			log.Println(err)
		}
		for _, trailer := range film.Trailers {
			_, err := repo.db.Exec(
				"INSERT INTO Trailers (KinopoiskID, Site, Title, URL) VALUES (?, ?, ?, ?)",
				trailer.KinopoiskID, trailer.Site, trailer.Title, trailer.URL,
			)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func (repo *FilmRepository) SaveUser(user models.User) {
	repo.db.Exec(
		"INSERT INTO Users (UserId, ChatId, UserName) VALUES (?, ?, ?)",
		user.UserID, user.ChatID, user.UserName,
	)
}

func (repo *FilmRepository) SaveFilmToUserPool(film models.FilmsPool) error {
	_, err := repo.db.Exec(
		"INSERT INTO UserPool (UserID, KinopoiskID, KinopoiskHDID, Year, Title, Description, PoserUrl) VALUES (?, ?, ?, ?, ?, ?, ?)",
		film.UserID, film.KinopoiskID, film.KinopoiskHDID, film.Year, film.Title, film.Description, film.PosterUrl,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *FilmRepository) DeleteFilmUserPool(filmId int, userId int64) error {
	_, err := repo.db.Exec("DELETE FROM UserPool WHERE KinopoiskID = ? AND UserID = ?", filmId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *FilmRepository) GetUserPool(userId int64, watch string) ([]models.FilmsPool, error) {
	rows, err := repo.db.Query("SELECT * FROM UserPool WHERE UserID = ? ", userId)
	if err != nil {
		return nil, err
	}
	films, err := scannerUserPoolFilm(rows)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func scannerUserPoolFilm(rows *sql.Rows) ([]models.FilmsPool, error) {
	var films []models.FilmsPool
	for rows.Next() {
		var film models.FilmsPool
		err := rows.Scan(&film.UserID, &film.KinopoiskID, &film.KinopoiskHDID, &film.Year, &film.Title, &film.Description, &film.PosterUrl, &film.Watch)
		if err != nil {
			return films, err
		}
		films = append(films, film)
	}
	return films, nil
}

func scannerFilm(row *sql.Row) (models.Film, error) {
	var film models.Film
	err := row.Scan(
		&film.KinopoiskID, &film.KinopoiskHDID, &film.TitleRu, &film.TitleOrig, &film.PosterUrl,
		&film.Description, &film.Country, &film.Year, &film.Rating,
		&film.Genre, &film.Posted,
	)
	if err != nil {
		return models.Film{}, err
	}
	return film, nil
}

func scannerTrailer(rows *sql.Rows) ([]models.Trailer, error) {
	var trailers []models.Trailer
	for rows.Next() {
		var trailer models.Trailer
		err := rows.Scan(&trailer.KinopoiskID, &trailer.URL, &trailer.Title, &trailer.Site)
		if err != nil {
			return nil, err
		}
		trailers = append(trailers, trailer)
	}

	return trailers, nil
}
