package db

import (
	"database/sql"
	"os"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteAdapter struct {
	db     *sql.DB
	dbName string
}

func NewSqliteAdapter(dbName string) (*SqliteAdapter, error) {
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	db.Exec(`
        CREATE TABLE invoices_requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        method VARCHAR(50) NOT NULL,
        url TEXT NOT NULL,
        remote_address VARCHAR(255),
        remote_port INT
    );`)

	db.Exec(`
        CREATE TABLE invoices_responses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        status_code INT NOT NULL
    );`)

	db.Exec(`
        CREATE TABLE invoices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        time TIMESTAMP NOT NULL,
        level iNT NOT NULL,
        pid INT NOT NULL,
        hostname VARCHAR(255) NOT NULL,
        response_time INT NOT NULL,
        message TEXT,
        request_id INT REFERENCES invoices_requests(id),
        response_id INT REFERENCES invoices_responses(id)
    );`)

	return &SqliteAdapter{db: db, dbName: dbName}, nil
}

func (m *SqliteAdapter) Close() error {
	return m.db.Close()
}

func (m *SqliteAdapter) Save(server string, invoice domain.Invoice) error {
	return nil
}

func (m *SqliteAdapter) Get(server string) ([]domain.Invoice, error) {
	return []domain.Invoice{}, nil
}

func (m *SqliteAdapter) Servers() ([]string, error) {
	return []string{}, nil
}

func (m *SqliteAdapter) Export() ([]domain.Invoice, error) {
	return []domain.Invoice{}, nil
}
