package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/serializers"
	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s       shortener.ShortenerService
	BaseURL string
	*echo.Echo
}

func NewHandler(s shortener.ShortenerService, bURL string) *handler {
	return &handler{
		s:       s,
		BaseURL: bURL,
		Echo:    echo.New(),
	}
}

func (h *handler) SetupRouting() {
	h.Echo.GET("/:url", h.GetURL)
	h.Echo.POST("/", h.PostURL)
	h.Echo.POST("/api/shorten", h.JSONPost)
}

func (h *handler) GetURL(c echo.Context) error {
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

func (h *handler) PostURL(c echo.Context) error {
	u := c.FormValue("url")
	_, err := url.Parse(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusInternalServerError))
	}

	sURL, err := h.s.Store(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.String(http.StatusCreated, fmt.Sprintf("%s/%s", h.BaseURL, sURL))
}

func (h *handler) JSONPost(c echo.Context) error {
	s, err := serializers.GetSerializer("json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	url, err := s.Decode(b)
	if err != nil {
		if errors.Is(err, json.ErrorEmptyURL) {
			return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	sURL, err := h.s.Store(url.URL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	jsonSURL := &model.ShortURL{ShortURL: fmt.Sprintf("%s/%s", h.BaseURL, sURL)}
	b, err = s.Encode(jsonSURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)
	_, err = c.Response().Write(b)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
