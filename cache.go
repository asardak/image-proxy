package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"time"
)

type Cache struct {
	root    string
	timeout time.Duration
}

func NewCache(path string, timeout time.Duration) *Cache {
	return &Cache{root: path, timeout: timeout}
}

func (c *Cache) Get(url string) (f io.ReadCloser, age time.Duration, err error) {
	file, err := os.Open(c.buildPath(url))
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}

	f = file
	info, err := file.Stat()
	if err != nil {
		return
	}

	age = time.Since(info.ModTime())
	if age > c.timeout {
		err = f.Close()
		f = nil
		return
	}

	return
}

func (c *Cache) NewWriter(url string) (io.WriteCloser, error) {
	return os.Create(c.buildPath(url))
}

func (c *Cache) Timeout() time.Duration {
	return c.timeout
}

func (c *Cache) buildPath(url string) string {
	return path.Join(c.root, hash(url)+".jpg")
}

func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
