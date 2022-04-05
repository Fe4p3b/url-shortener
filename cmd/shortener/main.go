package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	"github.com/Fe4p3b/url-shortener/internal/storage/pg"
	env "github.com/caarlos0/env/v6"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

type Config struct {
	Address         string `env:"SERVER_ADDRESS,required" envDefault:"0.0.0.0:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH,required" envDefault:"/tmp/url_shortener_storage" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN,required" envDefault:"postgres://postgres:12345@localhost:5432/shortener?sslmode=disable" json:"database_dsn"`
	Secret          string `env:"SECRET,required" envDefault:"x35k9f" json:"secret"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS,required" envDefault:"false" json:"enable_https"`
	Certfile        string `env:"CERTFILE" envDefault:"cert" json:"certfile_path"`
	CertKey         string `env:"PRIVATE_KEY" envDefault:"key" json:"certkey_path"`
	ConfigFile      string `env:"CONFIG" envDefault:"config/config.json"`
}

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

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

	if err = pg.CreateShortenerTable(); err != nil {
		log.Fatal(err)
	}

	s := shortener.NewShortener(pg, cfg.BaseURL)

	auth, err := auth.NewAuth([]byte(cfg.Secret), pg)
	if err != nil {
		log.Fatal(err)
	}
	authMiddleware := middleware.NewAuthMiddleware(auth)

	h := handlers.NewHandler(s)
	h.Router.Use(middleware.GZIPReaderMiddleware, middleware.GZIPWriterMiddleware, authMiddleware.Middleware)
	h.SetupAPIRouting()
	h.SetupProfiling()

	if cfg.EnableHTTPS {
		if err := createCert(); err != nil {
			log.Fatal(err)
		}

		if err := http.ListenAndServeTLS(cfg.Address, cfg.Certfile, cfg.CertKey, h.Router); err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return
	}

	if err := http.ListenAndServe(cfg.Address, h.Router); err != http.ErrServerClosed {
		log.Fatal(err)
	}
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
		enableHTTPS     bool
		configFile      string
	)

	flag.StringVar(&address, "a", "", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileStoragePath, "f", "", "Путь до файла с сокращёнными URL")
	flag.StringVar(&databaseDSN, "d", "", "Строка с адресом подключения к БД")
	flag.StringVar(&secret, "k", "", "Код для шифровки и дешифровки")
	flag.BoolVar(&enableHTTPS, "s", false, "Активация HTTPS")
	flag.StringVar(&configFile, "c", "", "Конфигурационный файл")
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

	if !enableHTTPS {
		cfg.EnableHTTPS = enableHTTPS
	}

	if configFile != "" {
		cfg.ConfigFile = configFile
	}

	if err := readJSONConfig(cfg); err != nil {
		return err
	}

	return nil
}

func readJSONConfig(cfg *Config) error {
	configFile, err := os.Open(cfg.ConfigFile)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(cfg); err != nil {
		return err
	}

	return nil
}

func createCert() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"KZ"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	f, err := os.OpenFile("cert", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = certPEM.WriteTo(f)
	if err != nil {
		return err
	}

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	f, err = os.OpenFile("key", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = privateKeyPEM.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}
