package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrBlankURL    = errors.New("url cannot be blank")
	ErrBlankWidth  = errors.New("width cannot be blank")
	ErrBlankHeight = errors.New("height cannot be blank")
)

type cacheInterface interface {
	Get(url string) (f io.ReadCloser, age time.Duration, err error)
	NewWriter(url string) (io.WriteCloser, error)
	Timeout() time.Duration
}

func GetImage(cache cacheInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, width, height, err := parseRequest(r)
		if err != nil {
			writeErr(w, err)
			return
		}

		file, expire, err := cache.Get(r.RequestURI)
		if err != nil {
			writeErr(w, err)
			return
		}

		if file != nil {
			defer file.Close()
			writeCacheHeaders(w, expire)
			io.Copy(w, file)
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			writeErr(w, err)
			return
		}
		defer resp.Body.Close()

		cacheWriter, err := cache.NewWriter(r.RequestURI)
		if err != nil {
			writeErr(w, err)
			return
		}
		defer cacheWriter.Close()

		multiWriter := io.MultiWriter(w, cacheWriter)

		writeCacheHeaders(w, cache.Timeout())
		err = Resize(resp.Body, width, height, multiWriter)
		if err != nil {
			writeErr(w, err)
			return
		}
	}
}

func parseRequest(r *http.Request) (url string, width int, height int, err error) {
	url = r.URL.Query().Get("url")
	if url == "" {
		err = ErrBlankURL
		return
	}

	if r.URL.Query().Get("width") == "" {
		err = ErrBlankWidth
		return
	}

	width, err = strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		return
	}

	if r.URL.Query().Get("height") == "" {
		err = ErrBlankHeight
		return
	}

	height, err = strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		return
	}

	return
}

func writeErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error: " + err.Error()))
}

func writeCacheHeaders(w http.ResponseWriter, expireTimeout time.Duration) {
	const (
		CacheControlHeader = "Cache-Control"
		ExpiresHeader      = "Expires"
	)

	w.Header().Add(CacheControlHeader, "public")
	w.Header().Add(ExpiresHeader, time.Now().Add(expireTimeout).Format(time.RFC1123Z))
}
