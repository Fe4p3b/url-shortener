package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	env "github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required" envDefault:"/tmp/url_shortener_storage"`
}

func setConfig(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	log.Printf("config from env: %v", cfg)

	var (
		address         string
		baseUrl         string
		fileStoragePath string
	)

	flag.StringVar(&address, "a", "localhost:8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseUrl, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileStoragePath, "f", "/tmp/url_shortener_storage", "Путь до файла с сокращёнными URL")
	flag.Parse()

	if address != "localhost:8080" {
		cfg.Address = address
	}

	if baseUrl != "http://localhost:8080" {
		cfg.BaseURL = baseUrl
	}

	if fileStoragePath != "/tmp/url_shortener_storage" {
		cfg.FileStoragePath = fileStoragePath
	}

	log.Printf("config from flags: %v", cfg)
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
	h.Router.Use(middleware.GZIPReaderMiddleware, middleware.GZIPWriterMiddleware)
	h.SetupRouting()

	if err := http.ListenAndServe(cfg.Address, h.Router); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
