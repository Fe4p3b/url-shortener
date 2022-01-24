package pg

import (
	"context"
	"database/sql"
	"log"
	"os"
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

func (p *pg) CreateShortenerTable() error {
	sql, err := os.ReadFile("./migrations/001_migration.sql")
	if err != nil {
		return err
	}
	log.Printf("migrations - %s", sql)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = p.db.ExecContext(ctx, string(sql))
	if err != nil {
		return err
	}

	return nil
}

func (p *pg) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := p.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (p *pg) Find(sURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `SELECT url FROM shortener.shortener WHERE short_url=$1`

	var (
		URL string
	)

	row := p.db.QueryRowContext(ctx, sql, sURL)

	if err := row.Scan(&URL); err != nil {
		log.Println(err)
		return "", err
	}

	return URL, nil
}

func (p *pg) Save(sURL string, URL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `INSERT INTO shortener.shortener VALUES($1, $2, $3)`

	if _, err := p.db.ExecContext(ctx, sql, sURL, URL, "asdf"); err != nil {
		log.Printf("err - %v", err)
		return err
	}

	return nil
}

func (p *pg) Close() error {
	return p.db.Close()
}
