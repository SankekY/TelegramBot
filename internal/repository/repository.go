package repository

import (
	"TelegramBot/internal/models"
	"database/sql"
	"fmt"
)

type FilmRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{db}
}

func (repo *FilmRepository) GetFilm() (models.Film, []models.Trailer, error) {
	film, err := scannerFilm(repo.db.QueryRow("SELECT * FROM Films LIMIT 1"))
	if err != nil {
		return models.Film{}, nil, fmt.Errorf("failed to retrieve film: %w", err)
	}

	rows, err := repo.db.Query("SELECT * FROM Trailers WHERE KinopoiskID = ?", film.KinopoiskID)
	if err != nil {
		return film, nil, fmt.Errorf("failed to retrieve trailers: %w", err)
	}
	defer rows.Close()
	trailers, err := scannerTrailer(rows)
	if err != nil {
		return film, nil, err
	}

	repo.db.Exec("update Films set Posted = 'true' WHERE KinopoiskID = ?", film.KinopoiskID)

	repo.db.QueryRow("SELECT COUNT(*) AS filmPosted FROM Films WHERE Posted = 'false'").Scan(&film.NonPosted)

	return film, trailers, nil
}

func (repo *FilmRepository) SaveFilmsAndTrailers(films []models.Film) error {
	for _, film := range films {
		film.Posted = "false"
		_, err := repo.db.Exec(
			"INSERT INTO Films (KinopoiskID, KinopoiskHDID, TitleRu, TitleOrig, PosterUrl, Description, Country, Year, Rating, Genre, Posted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			film.KinopoiskID, film.KinopoiskHDID, film.TitleRu, film.TitleOrig,
			film.PosterUrl, film.Description, film.Country,
			film.Year, film.Rating, film.Genre, film.Posted,
		)
		if err != nil {
			return err
		}
		for _, trailer := range film.Trailers {
			_, err := repo.db.Exec(
				"INSERT INTO Trailers (KinopoiskID, Site, Title, URL) VALUES (?, ?, ?, ?)",
				trailer.KinopoiskID, trailer.Site, trailer.Title, trailer.URL,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func scannerFilm(row *sql.Row) (models.Film, error) {
	var film models.Film
	err := row.Scan(
		&film.KinopoiskID, &film.KinopoiskHDID, &film.TitleRu, &film.TitleOrig,
		&film.PosterUrl, &film.Description, &film.Country, &film.Rating, &film.Year,
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
