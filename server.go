package main

import (
	"net/http"
)

func Server(cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()

		h := use(
			logRequests,
		)

		mux.HandleFunc("/image.jpg", h(GetImage(cache)))
		mux.HandleFunc("/", h(http.NotFound))

		mux.ServeHTTP(w, r)
	}
}
