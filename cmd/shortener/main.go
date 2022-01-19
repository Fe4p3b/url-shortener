package main

import (
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/server"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	env "github.com/caarlos0/env/v6"
)

const addr = "localhost:8080"

type Config struct {
	Address string `env:"SERVER_ADDRESS,required"`
	BaseURL string `env:"BASE_URL,required"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	m := memory.New(map[string]string{})
	s := shortener.New(m)
	h := handlers.New(s, cfg.BaseURL)
	h.SetupRouting()
	h.SetAddr(cfg.Address)

	server := server.New(h.Server)
	if err := server.Start(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
