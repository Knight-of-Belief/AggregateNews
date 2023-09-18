package api

import (
	"aggregatenews/pkg/store"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	db *store.DB
	r  *mux.Router
}

func New(db *store.DB) *API {
	api := API{
		db: db,
		r:  mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)
	nStr, ok := vars["n"]
	if !ok {
		http.Error(w, "Missing 'n' parameter", http.StatusBadRequest)
		return
	}
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	news, err := api.db.ReturnNews(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
