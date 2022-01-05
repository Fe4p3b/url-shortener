package handlers

import (
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/labstack/echo/v4"
)

type httpHandler struct {
	s shortener.ShortenerService
	*echo.Echo
}

func NewHttpHandler(s shortener.ShortenerService) *httpHandler {
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
		// return echo.NewHTTPError(http.StatusNotFound, err.Error())
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	// return c.String(http.StatusAccepted, url)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *httpHandler) EchoPost(c echo.Context) error {
	u := c.FormValue("url")
	uu, err := url.Parse(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if err == nil && uu.Host == "" && uu.Scheme == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL")
	}

	sUrl, err := h.s.Store(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.String(http.StatusCreated, sUrl)
}
