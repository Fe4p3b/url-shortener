package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/server"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	env "github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required" envDefault:"/tmp/url_shortener_storage"`
}

func setConfig(cfg *Config) error {
	if len(os.Args) > 1 {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "Адрес запуска HTTP-сервера")
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")
		flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/url_shortener_storage", "Путь до файла с сокращёнными URL")
		flag.Parse()
		return nil
	}

	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cfg := &Config{}
	if err := setConfig(cfg); err != nil {
		log.Fatal(err)
	}

	f, err := file.NewFile(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := shortener.NewShortener(f)
	h := handlers.NewHandler(s, cfg.BaseURL)
	h.SetupRouting()
	h.Use(middleware.GZIPmiddleware)

	server := server.NewServer(h.Server)
	if err := server.Start(cfg.Address); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
