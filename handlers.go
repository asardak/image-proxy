package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

var (
	ErrBlankURL    = errors.New("url cannot be blank")
	ErrBlankWidth  = errors.New("width cannot be blank")
	ErrBlankHeight = errors.New("height cannot be blank")
)

func (api *API) GetImage(w http.ResponseWriter, r *http.Request) {
	url, width, height, err := parseRequest(r)
	if err != nil {
		writeErr(w, err)
		return
	}

	_, _ = width, height

	file, err := api.cache.Get(url)
	defer file.Close()

	if file != nil {
		io.Copy(w, file)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		writeErr(w, err)
		return
	}
	defer resp.Body.Close()

	wc, err := api.cache.newWriter(url)
	if err != nil {
		writeErr(w, err)
		return
	}
	defer wc.Close()

	io.Copy(io.MultiWriter(w, wc), resp.Body)
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
