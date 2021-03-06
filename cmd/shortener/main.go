package main

import (
	"bytes"
	"context"
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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	grpcHandler "github.com/Fe4p3b/url-shortener/internal/handlers/grpc"
	pb "github.com/Fe4p3b/url-shortener/internal/handlers/grpc/proto"
	httpHandler "github.com/Fe4p3b/url-shortener/internal/handlers/http"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/file"
	"github.com/Fe4p3b/url-shortener/internal/storage/pg"
	env "github.com/caarlos0/env/v6"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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
	TrustedNetworks string `env:"TRUSTED_SUBNET" envDefault:"192.168.1.1" json:"trusted_subnet"`
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

	handlers := handlers.NewHandler(s)

	h := httpHandler.NewHandler(handlers)
	h.Router.Use(middleware.GZIPReaderMiddleware, middleware.GZIPWriterMiddleware, authMiddleware.Middleware)
	h.SetupAPIRouting()
	h.SetupProfiling()
	h.SetupInternalRouting(strings.Fields(cfg.TrustedNetworks))

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	shortenerServer := grpcHandler.NewShortenerServer(handlers)
	pb.RegisterShortenerServer(grpcServer, shortenerServer)

	srv := &http.Server{Addr: cfg.Address, Handler: h.Router}

	errgroup, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	errgroup.Go(func() error {
		return grpcServer.Serve(listen)
	})

	errgroup.Go(func() error {
		if cfg.EnableHTTPS {
			if err := createCert(); err != nil {
				return err
			}

			if err := srv.ListenAndServeTLS(cfg.Certfile, cfg.CertKey); err != http.ErrServerClosed {
				return err
			}

			return nil
		}

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	errgroup.Go(func() error {
		<-ctx.Done()
		log.Println("Shutting down server gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		grpcServer.GracefulStop()

		return srv.Shutdown(ctx)
	})

	if err := errgroup.Wait(); err != nil {
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
		trustedNetworks string
	)

	flag.StringVar(&address, "a", "", "?????????? ?????????????? HTTP-??????????????")
	flag.StringVar(&baseURL, "b", "", "?????????????? ?????????? ?????????????????????????????? ???????????????????????? URL")
	flag.StringVar(&fileStoragePath, "f", "", "???????? ???? ?????????? ?? ???????????????????????? URL")
	flag.StringVar(&databaseDSN, "d", "", "???????????? ?? ?????????????? ?????????????????????? ?? ????")
	flag.StringVar(&secret, "k", "", "?????? ?????? ???????????????? ?? ????????????????????")
	flag.BoolVar(&enableHTTPS, "s", false, "?????????????????? HTTPS")
	flag.StringVar(&configFile, "c", "", "???????????????????????????????? ????????")
	flag.StringVar(&trustedNetworks, "t", "", "IP-?????????????? ???????????????????? ??????????")
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

	if enableHTTPS {
		cfg.EnableHTTPS = enableHTTPS
	}

	if configFile != "" {
		cfg.ConfigFile = configFile
	}

	if trustedNetworks != "" {
		cfg.TrustedNetworks = trustedNetworks
	}

	if err := readJSONConfig(cfg); err != nil {
		return err
	}

	return nil
}

func readJSONConfig(cfg *Config) error {
	if cfg.ConfigFile != "" {
		configFile, err := os.Open(cfg.ConfigFile)
		if err != nil {
			return err
		}

		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(cfg); err != nil {
			return err
		}
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
