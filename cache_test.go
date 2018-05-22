package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

var cacheRoot = "/tmp/testdata"

func init() {
	if os.Getenv("IMG_PROXY_ROOT") != "" {
		cacheRoot = os.Getenv("IMG_PROXY_ROOT")
	}

	err := os.Mkdir(cacheRoot, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func TestCache_Get(t *testing.T) {
	const (
		imgURL  = "http://ya.ru/image1.jpg"
		imgData = "TestCache_Get data"
	)

	imgPath := path.Join(cacheRoot, hash(imgURL)) + ".jpg"
	file, err := os.Create(imgPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write([]byte(imgData))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	var f io.ReadCloser
	cache := NewCache(cacheRoot, time.Hour)
	f, _, err = cache.Get(imgURL)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	if string(data) != imgData {
		t.Fatalf("Expected data: %v, got: %v", imgData, string(data))
	}

	cache = NewCache(cacheRoot, time.Millisecond*500)
	f, _, err = cache.Get(imgURL)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)
	f, age, err := cache.Get(imgURL)
	if err != nil {
		t.Fatal(err)
	}

	if age < cache.timeout {
		t.Fatalf("Expected age is greater than cache.timeout %v", age)
	}

	if f != nil {
		t.Fatalf("Expected file is expired and equal to nil")
	}
}

func TestCache_NewWriter(t *testing.T) {
	const (
		imgURL  = "http://ya.ru/image2.jpg"
		imgData = "TestCache_NewWriter data"
	)

	cache := NewCache(cacheRoot, time.Hour)
	w, err := cache.NewWriter(imgURL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(imgData))
	if err != nil {
		t.Fatal(err)
	}

	f, _, err := cache.Get(imgURL)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != imgData {
		t.Fatalf("Expected cache data: %v, got: %v", imgData, string(data))
	}
}
