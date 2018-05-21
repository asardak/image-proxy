package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type Cache struct {
	root    string
	timeout time.Duration
	mu      sync.RWMutex
}

func NewCache(path string, timeout time.Duration) *Cache {
	return &Cache{root: path, timeout: timeout}
}

func (c *Cache) Get(url string) (*os.File, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f, err := os.Open(c.buildPath(url))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if time.Since(info.ModTime()) > c.timeout {
		return nil, f.Close()
	}

	return f, nil
}

func (c *Cache) newWriter(url string) (io.WriteCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return os.Create(c.buildPath(url))
}

func (c *Cache) buildPath(url string) string {
	return path.Join(c.root, hash(url)+".jpg")
}

func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
