package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type pg struct {
	db *sql.DB
}

var _ repositories.ShortenerRepository = &pg{}

func NewConnection(dsn string) (*pg, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &pg{db: conn}, nil
}

func (p *pg) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := p.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (p *pg) Find(string) (string, error) {
	return "", nil
}

func (p *pg) Save(string, string) error {
	return nil
}

func (p *pg) Close() error {
	return p.db.Close()
}
