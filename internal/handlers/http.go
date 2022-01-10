package handlers

import (
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/labstack/echo/v4"
)

type httpHandler struct {
	s shortener.ShortenerService
	*echo.Echo
}

func NewHTTPHandler(s shortener.ShortenerService) *httpHandler {
	h := &httpHandler{
		s:    s,
		Echo: echo.New(),
	}
	h.setupRouting()

	return h
}

func (h *httpHandler) setupRouting() {
	h.Echo.GET("/:url", h.EchoGet)
	h.Echo.POST("/", h.EchoPost)
}

func (h *httpHandler) EchoGet(c echo.Context) error {
	q := c.Param("url")
	if q == "" {
		return echo.NewHTTPError(http.StatusNotFound, "The query parameter is missing")
	}

	url, err := h.s.Find(q)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *httpHandler) EchoPost(c echo.Context) error {
	u := c.FormValue("url")

	sURL, err := h.s.Store(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.String(http.StatusCreated, sURL)
}
