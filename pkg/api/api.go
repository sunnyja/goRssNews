package api

import (
	"encoding/json"
	"goRssNews/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	db *storage.DB
	r 	*mux.Router
}


func New(db *storage.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a
}

//регистрация методов api в маршрутизаторе
func (api *API) endpoints() {
	//получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.news).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}


//получаем новости из базы
func (api *API) news(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, valErr := strconv.Atoi(s)
	if valErr != nil {
		http.Error(w, valErr.Error(), http.StatusUnsupportedMediaType)
	}
	//запрашиваем из БД n новостей
	news, err := api.db.News(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}


// Router возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.r
}