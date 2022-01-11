package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/labstack/echo/v4"
)

type httpHandler struct {
	s shortener.ShortenerService
	*echo.Echo
}

func NewHTTPHandler(s shortener.ShortenerService) *httpHandler {
	return &httpHandler{
		s:    s,
		Echo: echo.New(),
	}
}

func (h *httpHandler) SetAddr(addr string) {
	h.Server.Addr = addr
}

func (h *httpHandler) SetupRouting() {
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
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	uu, err := url.Parse(u)
	if err != nil || uu.Scheme == "" || uu.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL")
	}

	sURL, err := h.s.Store(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.String(http.StatusCreated, fmt.Sprintf("http://%s/%s", h.Server.Addr, sURL))
}
