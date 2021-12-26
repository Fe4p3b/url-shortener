package main

import (
	"log"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/handlers"
	"github.com/Fe4p3b/url-shortener/internal/server"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
)

func main() {
	m := memory.New()
	s := shortener.New(m)
	h := handlers.NewHttpHandler(s)
	_ = h
	server := server.New(":8080", nil)

	log.Fatal(server.Start())
}
