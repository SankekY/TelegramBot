package kinopoisk

import (
	"TelegramBot/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type KinopoiskAPI struct {
	token string
}

var indexTop int
var TOP_LIST = [13]string{
	"TOP_POPULAR_ALL", "TOP_POPULAR_MOVIES", "TOP_250_TV_SHOWS",
	"TOP_250_MOVIES", "VAMPIRE_THEME", "COMICS_THEME", "CLOSES_RELEASES", "FAMILY",
	"OSKAR_WINNERS_2021", "LOVE_THEME", "ZOMBIE_THEME", "CATASTROPHE_THEME", "POPULAR_SERIES",
}

func NewKinopoisk(token string) *KinopoiskAPI {
	return &KinopoiskAPI{token}
}

func (k *KinopoiskAPI) GetStackFilms() []models.Film {
	log.Println("GetStakFilms")
	var films []models.Film
	if indexTop >= 12 {
		indexTop = 0
	}
	FilmsStack := k.getReqStakFilms(TOP_LIST[indexTop])
	indexTop += 1
	log.Println("FilmsStack")
	for _, film := range FilmsStack {
		filmInfo, err := k.GetFullInfoByID(film.KinopoiskID)
		log.Println("Film info get ")
		if err != nil {
			log.Println(err)
		}

		for _, contry := range film.Countries {
			filmInfo.Country = contry.Country
			break
		}
		for _, genres := range film.Genres {
			filmInfo.Genre = genres.Genre
			break
		}
		filmInfo.Trailers = k.GetReqTrailers(film.KinopoiskID)
		films = append(films, filmInfo)
	}

	return films
}

func (k *KinopoiskAPI) getReqStakFilms(topic string) []Item {
	client := &http.Client{}
	url := fmt.Sprintf("https://kinopoiskapiunofficial.tech/api/v2.2/films/collections?type=%s&page=1", topic)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-KEY", k.token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var r Response
	if err := json.Unmarshal(respByte, &r); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	log.Println("getReqStak")
	return r.Items
}

func (k *KinopoiskAPI) GetReqTrailers(id int) []models.Trailer {
	client := &http.Client{}
	url := fmt.Sprintf("https://kinopoiskapiunofficial.tech/api/v2.2/films/%d/videos", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil
	}
	req.Header.Set("X-API-KEY", k.token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	log.Printf("Method: %s [%d]", resp.Request.Method, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil
	}

	var body ResponseTraller
	if err := json.Unmarshal(data, &body); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return nil
	}

	var trailers []models.Trailer
	for _, item := range body.Items {
		if item.Site != "UNKNOWN" {
			trailers = append(trailers, models.Trailer{
				KinopoiskID: id,
				URL:         item.Url,
				Title:       item.Name,
				Site:        item.Site,
			})
		}
	}
	return trailers
}

func (k *KinopoiskAPI) GetFullInfoByID(id int) (models.Film, error) {
	url := fmt.Sprintf("https://kinopoiskapiunofficial.tech/api/v2.2/films/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return models.Film{}, err
	}
	req.Header.Set("X-API-KEY", k.token)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return models.Film{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return models.Film{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return models.Film{}, err
	}

	var body models.Film
	if err := json.Unmarshal(data, &body); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return models.Film{}, err
	}
	log.Printf("Fetched full info for film ID: %d", id)
	return body, nil
}
