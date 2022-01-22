package middleware

import (
	"compress/gzip"
	"io"
	"log"
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

func GZIPmiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().Header.Get(echo.HeaderAcceptEncoding))
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
