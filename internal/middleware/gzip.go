package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZIPWriterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "gzip") {
			return next(c)
		}

		gz, err := gzip.NewWriterLevel(c.Response().Writer, gzip.BestSpeed)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		defer gz.Close()

		c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
		c.Response().Writer = gzipWriter{ResponseWriter: c.Response().Writer, Writer: gz}
		return next(c)
	}
}

func GZIPReaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderContentEncoding) == "gzip" {
			gz, err := gzip.NewReader(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			}
			defer gz.Close()

			c.Request().Body = gz
		}

		return next(c)
	}
}
