package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	"github.com/Fe4p3b/url-shortener/internal/storage/pg"
	env "github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required" envDefault:"/tmp/url_shortener_storage"`
	DatabaseDSN     string `env:"DATABASE_DSN,required" envDefault:"postgres://postgres:12345@localhost:5432/shortener?sslmode=disable"`
}

func setConfig(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	var (
		address         string
		baseURL         string
		fileStoragePath string
		databaseDSN     string
	)

	flag.StringVar(&address, "a", "localhost:8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileStoragePath, "f", "/tmp/url_shortener_storage", "Путь до файла с сокращёнными URL")
	flag.StringVar(&databaseDSN, "d", "postgres://postgres:12345@localhost:5432/shortener?sslmode=disable", "Строка с адресом подключения к БД")
	flag.Parse()

	if address != "localhost:8080" {
		cfg.Address = address
	}

	if baseURL != "http://localhost:8080" {
		cfg.BaseURL = baseURL
	}

	if fileStoragePath != "/tmp/url_shortener_storage" {
		cfg.FileStoragePath = fileStoragePath
	}

	if databaseDSN != "postgres://postgres:12345@localhost:5432/shortener?sslmode=disable" {
		cfg.DatabaseDSN = databaseDSN
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

	pg, err := pg.NewConnection(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	if err := pg.CreateShortenerTable(); err != nil {
		log.Fatal(err)
	}

	s := shortener.NewShortener(pg, cfg.BaseURL)

	h := handlers.NewHandler(s)
	h.Router.Use(middleware.GZIPReaderMiddleware, middleware.GZIPWriterMiddleware, middleware.Auth)
	h.SetupRouting()

	if err := http.ListenAndServe(cfg.Address, h.Router); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
