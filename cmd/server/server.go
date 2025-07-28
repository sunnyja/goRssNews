package main

import (
	"encoding/json"
	"goRssNews/pkg/api"
	"goRssNews/pkg/rss"
	"goRssNews/pkg/storage"
	"log"
	"net/http"
	"os"
	"time"
)

// конфигурация приложения
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	db, err := storage.New("postgres://postgres:postgres@localhost:5432/goNews?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)
	conf, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config config
	err = json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatal(err)
	}

	// запуск парсинга новостей в отдельном потоке
	// для каждой ссылки
	chPosts := make(chan []storage.News)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, chPosts, chErrs, config.Period)
	}

	// запись потока новостей в БД
	go func() {
		for posts := range chPosts {
			db.AddNews(posts)
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}


// Асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, posts chan<- []storage.News, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}