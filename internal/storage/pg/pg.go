// Package pg implements postgres storage for
// shortener service.
package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// pg contains database conneciton, buffer for bulk
// addition and buffer to delete URLs.
type pg struct {
	// db is a database connection
	db *sql.DB

	// buffer for bulk addition
	buffer []repositories.URL

	// deleteBuffer for URLs deletion
	deleteBuffer chan repositories.URL
}

var _ repositories.ShortenerRepository = &pg{}
var _ repositories.AuthRepository = &pg{}

func NewConnection(dsn string) (*pg, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &pg{db: conn, buffer: make([]repositories.URL, 0, 1000), deleteBuffer: make(chan repositories.URL, 1)}, nil
}

// CreateShortenerTable creates required fields
// in database.
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

// Ping implements repositories.ShortenerRepository Ping method.
func (p *pg) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := p.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// Find implements repositories.ShortenerRepository Find method.
func (p *pg) Find(sURL string) (*repositories.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `SELECT original_url, is_deleted FROM shortener.shortener WHERE short_url=$1`

	URL := &repositories.URL{}

	row := p.db.QueryRowContext(ctx, sql, sURL)

	if err := row.Scan(&URL.URL, &URL.IsDeleted); err != nil {
		return nil, err
	}

	return URL, nil
}

// Save implements repositories.ShortenerRepository Save method.
func (p *pg) Save(url *models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `INSERT INTO shortener.shortener(short_url, original_url, user_id) VALUES($1, $2, $3)`

	_, err := p.db.ExecContext(ctx, sql, url.ShortURL, url.URL, url.UserID)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		sql := `SELECT short_url FROM shortener.shortener WHERE original_url=$1`
		row := p.db.QueryRowContext(ctx, sql, url.URL)
		if err = row.Scan(&url.ShortURL); err != nil {
			return err
		}
	}

	return err
}

// GetUserURLs implements repositories.ShortenerRepository GetUserURLs method.
func (p *pg) GetUserURLs(user string, baseURL string) (URLs []repositories.URL, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `SELECT short_url, original_url FROM shortener.shortener WHERE is_deleted=false and user_id=$1`

	rows, err := p.db.QueryContext(ctx, sql, user)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var URL repositories.URL
		if err := rows.Scan(&URL.ShortURL, &URL.URL); err != nil {
			return nil, err
		}

		URL.ShortURL = fmt.Sprintf("%s/%s", baseURL, URL.ShortURL)
		URLs = append(URLs, URL)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return
}

// AddURLBuffer implements repositories.ShortenerRepository AddURLBuffer method.
func (p *pg) AddURLBuffer(u repositories.URL) error {
	p.buffer = append(p.buffer, u)

	if len(p.buffer) == cap(p.buffer) {
		if err := p.Flush(); err != nil {
			return err
		}
	}

	return nil
}

// Flush implements repositories.ShortenerRepository Flush method.
func (p *pg) Flush() error {
	if len(p.buffer) == 0 {
		return nil
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO shortener.shortener(correlation_id, short_url, original_url, user_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	for _, v := range p.buffer {
		if _, err := stmt.Exec(v.CorrelationID, v.ShortURL, v.URL, v.UserID); err != nil {
			if err = tx.Rollback(); err != nil {
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

// AddURLToDelete implements repositories.ShortenerRepository AddURLToDelete method.
func (p *pg) AddURLToDelete(u repositories.URL) {
	p.deleteBuffer <- u
}

// FlushToDelete implements repositories.ShortenerRepository FlushToDelete method.
func (p *pg) FlushToDelete() error {
	if len(p.deleteBuffer) == 0 {
		return nil
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE shortener.shortener SET is_deleted=true WHERE short_url=$1 and user_id=$2")
	if err != nil {
		return err
	}

	for {
		select {
		case v := <-p.deleteBuffer:
			if _, err := stmt.Exec(v.ShortURL, v.UserID); err != nil {
				if err = tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		default:
			if err := tx.Commit(); err != nil {
				return err
			}
			return nil
		}
	}
}

// Close closes database connection.
func (p *pg) Close() error {
	return p.db.Close()
}

// CreateUser implements repositories.AuthRepository CreateUser method.
func (p *pg) CreateUser() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `INSERT INTO shortener.users VALUES(default) RETURNING id`

	row := p.db.QueryRowContext(ctx, sql)

	var uuid string

	if err := row.Scan(&uuid); err != nil {
		return "", err
	}

	return uuid, nil
}

// VerifyUser implements repositories.AuthRepository VerifyUser method.
func (p *pg) VerifyUser(user string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sql := `SELECT id FROM shortener.users WHERE id=$1`

	row := p.db.QueryRowContext(ctx, sql, user)
	var uuid string

	if err := row.Scan(&uuid); err != nil {
		return err
	}

	return nil
}

func (p *pg) GetStats() (*models.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := &models.Stats{}

	sql := `SELECT COUNT(short_url) FROM shortener.shortener`
	row := p.db.QueryRowContext(ctx, sql)
	if err := row.Scan(&stats.URLs); err != nil {
		return nil, err
	}

	sql = `SELECT COUNT(*) FROM shortener.users`
	row = p.db.QueryRowContext(ctx, sql)
	if err := row.Scan(&stats.Users); err != nil {
		return nil, err
	}

	return stats, nil
}
