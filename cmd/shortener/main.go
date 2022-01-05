package main

import (
	"log"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/server"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
)

func main() {
	m := memory.New(map[string]string{})
	s := shortener.New(m)
	h := handlers.NewHttpHandler(s)

	e := h.InitEchoHandler()

	server := server.New(":8080", e)
	if err := server.Start(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
