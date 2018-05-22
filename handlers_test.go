package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type cacheMock struct {
	GetFunc       func(url string) (f io.ReadCloser, age time.Duration, err error)
	NewWriterFunc func(url string) (io.WriteCloser, error)
	TimeoutFunc   func() time.Duration
}

func (m *cacheMock) Get(url string) (f io.ReadCloser, age time.Duration, err error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}

	return
}

func (m *cacheMock) NewWriter(url string) (io.WriteCloser, error) {
	if m.NewWriterFunc != nil {
		return m.NewWriterFunc(url)
	}

	return nil, nil
}

func (m *cacheMock) Timeout() time.Duration {
	if m.TimeoutFunc != nil {
		return m.TimeoutFunc()
	}

	return 0
}

type closeBuffer bytes.Buffer

func (b *closeBuffer) Close() error                      { return nil }
func (b *closeBuffer) Read(p []byte) (n int, err error)  { return (*bytes.Buffer)(b).Read(p) }
func (b *closeBuffer) Write(p []byte) (n int, err error) { return (*bytes.Buffer)(b).Write(p) }

func TestGetImage(t *testing.T) {
	cache := &cacheMock{}

	t.Run("With blank params", func(t *testing.T) {
		rec := httptest.NewRecorder()
		q := httptest.NewRequest("GET", "http://localhost:8080/image.jpg", nil)
		GetImage(cache)(rec, q)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("Expected status: %v, got: %v", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("With valid params", func(t *testing.T) {
		rec := httptest.NewRecorder()
		q := httptest.NewRequest("GET", "http://localhost:8080/image.jpg?width=100&height=100&url=http://ya.ru/image.jpg", nil)

		const testData = "TestGetImage_data"

		cache.GetFunc = func(url string) (f io.ReadCloser, age time.Duration, err error) {
			return (*closeBuffer)(bytes.NewBuffer([]byte(testData))), time.Second, nil
		}
		cache.TimeoutFunc = func() time.Duration { return time.Hour }

		GetImage(cache)(rec, q)

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status: %v, got: %v", http.StatusInternalServerError, rec.Code)
		}

		if string(rec.Body.Bytes()) != testData {
			t.Fatalf("Expected body: %v, got: %v", testData, string(rec.Body.Bytes()))
		}
	})

	t.Run("With cache miss", func(t *testing.T) {
		testData, err := base64.StdEncoding.DecodeString(srcImgBase64)
		if err != nil {
			t.Fatal(err)
		}

		resData, err := base64.StdEncoding.DecodeString(resImgBase64)
		if err != nil {
			t.Fatal(err)
		}

		cacheBuf := (*closeBuffer)(bytes.NewBuffer([]byte{}))
		cache.GetFunc = func(url string) (f io.ReadCloser, age time.Duration, err error) {
			return nil, 0, nil
		}
		cache.TimeoutFunc = func() time.Duration { return time.Hour }
		cache.NewWriterFunc = func(url string) (io.WriteCloser, error) {
			return cacheBuf, nil
		}

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testData))
		}))

		rec := httptest.NewRecorder()
		q := httptest.NewRequest("GET", "http://localhost:8080/image.jpg?width=50&height=50&url="+s.URL, nil)

		GetImage(cache)(rec, q)

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status: %v, got: %v", http.StatusOK, rec.Code)
		}

		if !reflect.DeepEqual(rec.Body.Bytes(), resData) {
			t.Fatalf("Expected body: %v, got: %v", resImgBase64, string(rec.Body.Bytes()))
		}

		if !reflect.DeepEqual((*bytes.Buffer)(cacheBuf).Bytes(), resData) {
			t.Fatalf("Expected cached data: %v, got: %v", resImgBase64, string((*bytes.Buffer)(cacheBuf).Bytes()))
		}
	})
}
