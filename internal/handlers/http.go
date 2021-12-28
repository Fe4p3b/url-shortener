package handlers

import (
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
)

type httpHandler struct {
	s shortener.ShortenerService
}

func NewHttpHandler(s shortener.ShortenerService) *httpHandler {
	h := &httpHandler{
		s: s,
	}
	http.HandleFunc("/", h.handler)
	return h
}

func (h *httpHandler) handler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(rw, r)
	}

	if r.Method == http.MethodPost {
		h.post(rw, r)
	}

	http.Error(rw, "", http.StatusMethodNotAllowed)
}

func (h *httpHandler) get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("url")
	if q == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}

	url, err := h.s.Find(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *httpHandler) post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u := r.FormValue("url")
	uu, err := url.Parse(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil && uu.Host == "" && uu.Scheme == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	sUrl, err := h.s.Store(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(sUrl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
