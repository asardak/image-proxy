package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	f := parseFlags()

	log.Fatalln(
		http.ListenAndServe(f.addr, Server(
			NewCache(f.cachePath, f.cacheTime),
		)),
	)
}

type flags struct {
	version   bool
	help      bool
	cachePath string
	cacheTime time.Duration
	addr      string
}

func parseFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.version, "v", false, "Show version")
	flag.BoolVar(&f.help, "h", false, "Help")
	flag.StringVar(&f.addr, "addr", ":8080", "Listen address")
	flag.StringVar(&f.cachePath, "cache-path", "/tmp/image-proxy", "Path to cache folder")
	flag.DurationVar(&f.cacheTime, "cache-time", time.Hour, "Cache timeout")
	flag.Parse()

	if f.version {
		fmt.Println("Image-proxy version:", Version)
		os.Exit(0)
	}

	if f.help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	return f
}
