package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	"github.com/Fe4p3b/url-shortener/internal/storage/pg"
	env "github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"SERVER_ADDRESS,required" envDefault:"0.0.0.0:8080"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required" envDefault:"/tmp/url_shortener_storage"`
	DatabaseDSN     string `env:"DATABASE_DSN,required" envDefault:"postgres://postgres:12345@localhost:5432/shortener?sslmode=disable"`
	Secret          string `env:"SECRET,required" envDefault:"x35k9f"`
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
		secret          string
	)

	flag.StringVar(&address, "a", "", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileStoragePath, "f", "", "Путь до файла с сокращёнными URL")
	flag.StringVar(&databaseDSN, "d", "", "Строка с адресом подключения к БД")
	flag.StringVar(&secret, "s", "", "Код для шифровки и дешифровки")
	flag.Parse()

	if address != "" {
		cfg.Address = address
	}

	if baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if fileStoragePath != "" {
		cfg.FileStoragePath = fileStoragePath
	}

	if databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}

	if secret != "" {
		cfg.Secret = secret
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

	auth := auth.NewAuth([]byte(cfg.Secret), pg)
	authMiddleware := middleware.NewAuthMiddleware(auth)

	h := handlers.NewHandler(s)
	h.Router.Use(middleware.GZIPReaderMiddleware, middleware.GZIPWriterMiddleware, authMiddleware.Middleware)
	h.SetupRouting()

	if err := http.ListenAndServe(cfg.Address, h.Router); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
