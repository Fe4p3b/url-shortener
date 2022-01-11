package main

import (
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/server"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
)

const addr = "localhost:8080"

func main() {

	m := memory.New(map[string]string{})
	s := shortener.New(m)
	h := handlers.NewHTTPHandler(s)
	h.SetupRouting()
	h.SetAddr(addr)

	server := server.New(h.Server)
	if err := server.Start(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
