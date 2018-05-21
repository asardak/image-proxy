package main

import (
	"net/http"
)

type API struct {
	cache *Cache
}

func NewAPI(cache *Cache) *API {
	return &API{
		cache: cache,
	}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	h := use(
		logRequests,
	)

	mux.HandleFunc("/image.jpg", h(api.GetImage))
	mux.HandleFunc("/", h(http.NotFound))

	mux.ServeHTTP(w, r)
}
