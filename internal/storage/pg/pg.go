package pg

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type pg struct {
	db     *sql.DB
	buffer []repositories.URL
}

var _ repositories.ShortenerRepository = &pg{}

func NewConnection(dsn string) (*pg, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &pg{db: conn, buffer: make([]repositories.URL, 0, 1000)}, nil
}

func (p *pg) CreateShortenerTable() error {
	sql, err := os.ReadFile("./migrations/001_migration.sql")
	if err != nil {
		return err
	}

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

	sql := `SELECT original_url FROM shortener.shortener WHERE short_url=$1`

	var (
		URL string
	)

	row := p.db.QueryRowContext(ctx, sql, sURL)

	if err := row.Scan(&URL); err != nil {
		return "", err
	}

	return URL, nil
}

func (p *pg) Save(sURL *string, URL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `INSERT INTO shortener.shortener(short_url, original_url, user_id) VALUES($1, $2, $3)`

	_, err := p.db.ExecContext(ctx, sql, sURL, URL, "asdf")
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		sql := `SELECT short_url FROM shortener.shortener WHERE original_url=$1`
		log.Printf("short_url before: %v", *sURL)
		row := p.db.QueryRowContext(ctx, sql, URL)
		if err := row.Scan(sURL); err != nil {
			return err
		}
		log.Printf("short_url after: %v", *sURL)
	}

	return err
}

func (p *pg) AddURLBuffer(u repositories.URL) error {
	p.buffer = append(p.buffer, u)

	if cap(p.buffer) == len(p.buffer) {
		if err := p.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func (p *pg) Flush() error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO shortener.shortener(correlation_id, short_url, original_url, user_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	for _, v := range p.buffer {
		if _, err := stmt.Exec(v.CorrelationID, v.ShortURL, v.URL, ""); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	p.buffer = p.buffer[:0]
	return nil
}

func (p *pg) Close() error {
	return p.db.Close()
}
