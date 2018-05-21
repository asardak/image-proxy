package main

import (
	"log"
	"net/http"
	"time"
)

type middleware func(next http.HandlerFunc) http.HandlerFunc

func use(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func logRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s\t%s\tlatency: %v", r.Method, r.RequestURI, time.Since(start))
	}
}
