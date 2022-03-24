// Package pg implements postgres storage for
// shortener service.
package pg

import (
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_pg_Ping(t *testing.T) {
	db, _ := NewMock()
	defer db.Close()
	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:           db,
				buffer:       make([]repositories.URL, 0),
				deleteBuffer: make(chan repositories.URL),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				buffer:       tt.fields.buffer,
				deleteBuffer: tt.fields.deleteBuffer,
			}
			err := p.Ping()
			assert.NoError(t, err)
		})
	}
}

func Test_pg_Find(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	type args struct {
		sURL  string
		query string
		URL   repositories.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *repositories.URL
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:           db,
				buffer:       make([]repositories.URL, 0),
				deleteBuffer: make(chan repositories.URL),
			},
			args: args{
				sURL:  "asdf",
				query: "SELECT original_url, is_deleted FROM shortener.shortener WHERE short_url=$1",
				URL: repositories.URL{
					URL:       "http://google.com",
					IsDeleted: false,
				},
			},
			want: &repositories.URL{
				URL:       "http://google.com",
				IsDeleted: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				buffer:       tt.fields.buffer,
				deleteBuffer: tt.fields.deleteBuffer,
			}
			rows := sqlmock.NewRows([]string{"original_url", "is_deleted"}).
				AddRow(tt.args.URL.URL, tt.args.URL.IsDeleted)
			mock.ExpectQuery(regexp.QuoteMeta(tt.args.query)).WithArgs(tt.args.sURL).WillReturnRows(rows)

			got, err := p.Find(tt.args.sURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_pg_Store(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	type args struct {
		sURL  string
		query string
		URL   models.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:           db,
				buffer:       make([]repositories.URL, 0),
				deleteBuffer: make(chan repositories.URL),
			},
			args: args{
				query: "INSERT INTO shortener.shortener(short_url, original_url, user_id) VALUES($1, $2, $3)",
				URL: models.URL{
					URL:      "http://google.com",
					UserID:   "1234",
					ShortURL: "asdf",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				buffer:       tt.fields.buffer,
				deleteBuffer: tt.fields.deleteBuffer,
			}

			prep := mock.ExpectExec(regexp.QuoteMeta(tt.args.query))
			prep.WithArgs(tt.args.URL.ShortURL, tt.args.URL.URL, tt.args.URL.UserID).WillReturnResult(sqlmock.NewResult(0, 1))

			err := p.Save(&tt.args.URL)
			assert.NoError(t, err)
		})
	}
}

func Test_pg_CreateUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type args struct {
		query string
		id    string
	}
	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:           db,
				buffer:       make([]repositories.URL, 0),
				deleteBuffer: make(chan repositories.URL),
			},
			args: args{
				query: "INSERT INTO shortener.users VALUES(default) RETURNING id",
				id:    "asdf",
			},
			want: "asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				buffer:       tt.fields.buffer,
				deleteBuffer: tt.fields.deleteBuffer,
			}

			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(tt.args.id)

			prep := mock.ExpectQuery(regexp.QuoteMeta(tt.args.query))
			prep.WillReturnRows(rows)

			got, err := p.CreateUser()
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func Test_pg_VerifyUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	type args struct {
		user  string
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				db:           db,
				buffer:       make([]repositories.URL, 0),
				deleteBuffer: make(chan repositories.URL),
			},
			args: args{
				query: "SELECT id FROM shortener.users WHERE id=$1",
				user:  "asdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				buffer:       tt.fields.buffer,
				deleteBuffer: tt.fields.deleteBuffer,
			}

			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(tt.args.user)
			mock.ExpectQuery(regexp.QuoteMeta(tt.args.query)).WithArgs(tt.args.user).WillReturnRows(rows)

			err := p.VerifyUser(tt.args.user)
			assert.NoError(t, err)
		})
	}
}

func Test_pg_AddURLToDelete(t *testing.T) {
	type fields struct {
		deleteBuffer chan repositories.URL
	}
	type args struct {
		u repositories.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   repositories.URL
		args   args
	}{
		{
			fields: fields{
				deleteBuffer: make(chan repositories.URL),
			},
			args: args{
				u: repositories.URL{
					URL:      "http://google.com",
					ShortURL: "asdf",
				},
			},
			want: repositories.URL{
				URL:      "http://google.com",
				ShortURL: "asdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				deleteBuffer: tt.fields.deleteBuffer,
			}
			go func() {
				p.AddURLToDelete(tt.args.u)
			}()

			got := <-p.deleteBuffer
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_pg_AddURLBuffer(t *testing.T) {
	db, _ := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	type args struct {
		u repositories.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []repositories.URL
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:     db,
				buffer: make([]repositories.URL, 0, 10),
			},
			args: args{
				u: repositories.URL{
					URL:      "http://google.com",
					ShortURL: "asdf",
				},
			},
			want: []repositories.URL{
				{
					URL:      "http://google.com",
					ShortURL: "asdf",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:     tt.fields.db,
				buffer: tt.fields.buffer,
			}
			err := p.AddURLBuffer(tt.args.u)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, p.buffer)
		})
	}
}

func Test_pg_Flush(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}

	tests := []struct {
		name    string
		fields  fields
		query   string
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db: db,
				buffer: []repositories.URL{
					{
						CorrelationID: "1",
						URL:           "http://google.com",
						ShortURL:      "asdf",
						UserID:        "1",
					},
				},
			},
			query:   "INSERT INTO shortener.shortener(correlation_id, short_url, original_url, user_id) VALUES($1, $2, $3, $4)",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:     tt.fields.db,
				buffer: tt.fields.buffer,
			}
			mock.ExpectBegin()
			prep := mock.ExpectPrepare(regexp.QuoteMeta(tt.query))
			for _, arg := range p.buffer {

				prep.ExpectExec().WithArgs(arg.CorrelationID, arg.ShortURL, arg.URL, arg.UserID).WillReturnResult(sqlmock.NewResult(0, 1))
			}
			mock.ExpectCommit()

			err := p.Flush()
			assert.NoError(t, err)
		})
	}
}

func Test_pg_GetUserURLs(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		buffer       []repositories.URL
		deleteBuffer chan repositories.URL
	}
	type args struct {
		user    string
		baseURL string
		query   string
		URL     repositories.URL
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantURLs []repositories.URL
		wantErr  bool
	}{
		{
			fields: fields{
				db: db,
			},
			args: args{
				user:    "asdf",
				baseURL: "localhost:8080",
				query:   "SELECT short_url, original_url FROM shortener.shortener WHERE is_deleted=false and user_id=$1",
				URL: repositories.URL{
					ShortURL: "qwer",
					URL:      "http://google.com",
				},
			},
			wantURLs: []repositories.URL{
				{
					ShortURL: "localhost:8080/qwer",
					URL:      "http://google.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db: tt.fields.db,
			}

			rows := sqlmock.NewRows([]string{"short_url", "original_url"}).
				AddRow(tt.args.URL.ShortURL, tt.args.URL.URL)
			mock.ExpectQuery(regexp.QuoteMeta(tt.args.query)).WithArgs(tt.args.user).WillReturnRows(rows)

			gotURLs, err := p.GetUserURLs(tt.args.user, tt.args.baseURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantURLs, gotURLs)
		})
	}
}

func Test_pg_FlushToDelete(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	type fields struct {
		db           *sql.DB
		deleteBuffer chan repositories.URL
	}
	type args struct {
		URL   repositories.URL
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		query   string
		args    args
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				db:           db,
				deleteBuffer: make(chan repositories.URL, 1),
			},
			args: args{
				URL: repositories.URL{
					CorrelationID: "1",
					URL:           "http://google.com",
					ShortURL:      "asdf",
					UserID:        "1",
				},
				query: "UPDATE shortener.shortener SET is_deleted=true WHERE short_url=$1 and user_id=$2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pg{
				db:           tt.fields.db,
				deleteBuffer: tt.fields.deleteBuffer,
			}

			go func() {
				p.deleteBuffer <- tt.args.URL
			}()

			time.Sleep(1 * time.Second)
			mock.ExpectBegin()
			prep := mock.ExpectPrepare(regexp.QuoteMeta(tt.query))
			prep.ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			err := p.FlushToDelete()
			assert.NoError(t, err)
		})
	}
}
