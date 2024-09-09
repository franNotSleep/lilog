package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteAdapter struct {
	db     *sql.DB
	dbName string
}

const (
	INVOICE_TABLE = "invoices"
)

const (
	YELLOW = "\033[33m"
	GREEN  = "\033[32m"
	RESET  = "\033[0m"
)

func NewSqliteAdapter(dbName string) (*SqliteAdapter, error) {
	err := os.Remove(dbName)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	log.Println(YELLOW + "Creating tables...⏳ " + RESET)
	_, err = db.Exec(`
        CREATE TABLE invoices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        server VARCHAR(255) NOT NULL,
        time TIMESTAMP NOT NULL,
        level iNT NOT NULL,
        hostname VARCHAR(255) NOT NULL,
        response_time INT NOT NULL,
        message TEXT,
        status_code INT NOT NULL,
        method VARCHAR(50) NOT NULL,
        url TEXT NOT NULL,
        remote_address VARCHAR(255),
        remote_port INT
    );`, INVOICE_TABLE)

	if err != nil {
		return nil, err
	}

	log.Println(GREEN + "Tables have been created successfully✅" + RESET)
	return &SqliteAdapter{db: db, dbName: dbName}, nil
}

func (m *SqliteAdapter) Close() error {
	return m.db.Close()
}

func (m *SqliteAdapter) Save(server string, invoice domain.Invoice) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := m.db.Prepare("INSERT INTO " + INVOICE_TABLE + " (server, time, level, hostname, response_time, message, status_code, method, url, remote_address, remote_port) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(server, invoice.Time, invoice.Level, invoice.Hostname, invoice.ResponseTime, invoice.Message, invoice.InvoiceResponse.StatusCode, invoice.InvoiceRequest.Method, invoice.InvoiceRequest.URL, invoice.InvoiceRequest.RemoteAddress, invoice.InvoiceRequest.RemotePort)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *SqliteAdapter) Get(server string) ([]domain.Invoice, error) {
	rows, err := m.db.Query("SELECT * FROM "+INVOICE_TABLE+" WHERE server=?", server)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []domain.Invoice{}

	for rows.Next() {
		invoice := domain.Invoice{
			InvoiceRequest:  domain.InvoiceRequest{},
			InvoiceResponse: domain.InvoiceResponse{},
		}

		if err := rows.Scan(&invoice.ID, &invoice.Server, &invoice.Time, &invoice.Level, &invoice.Hostname, &invoice.ResponseTime, &invoice.Message, &invoice.InvoiceResponse.StatusCode, &invoice.InvoiceRequest.Method, &invoice.InvoiceRequest.URL, &invoice.InvoiceRequest.RemoteAddress, &invoice.InvoiceRequest.RemotePort); err != nil {
			return invoices, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return invoices, err
	}

	return invoices, nil
}

func (m *SqliteAdapter) Servers() ([]string, error) {
	return []string{}, nil
}

func (m *SqliteAdapter) Export() ([]domain.Invoice, error) {
	rows, err := m.db.Query("SELECT * FROM " + INVOICE_TABLE + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []domain.Invoice{}

	for rows.Next() {
		invoice := domain.Invoice{
			InvoiceRequest:  domain.InvoiceRequest{},
			InvoiceResponse: domain.InvoiceResponse{},
		}

		if err := rows.Scan(&invoice.ID, &invoice.Server, &invoice.Time, &invoice.Level, &invoice.Hostname, &invoice.ResponseTime, &invoice.Message, &invoice.InvoiceResponse.StatusCode, &invoice.InvoiceRequest.Method, &invoice.InvoiceRequest.URL, &invoice.InvoiceRequest.RemoteAddress, &invoice.InvoiceRequest.RemotePort); err != nil {
			return invoices, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return invoices, err
	}

	return invoices, nil
}

func (m *SqliteAdapter) GetBetween(from time.Time, until time.Time) ([]domain.Invoice, error) {
	rows, err := m.db.Query("SELECT * FROM "+INVOICE_TABLE+"WHERE time BETWEEN ? AND ?;", from, until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []domain.Invoice{}

	for rows.Next() {
		invoice := domain.Invoice{
			InvoiceRequest:  domain.InvoiceRequest{},
			InvoiceResponse: domain.InvoiceResponse{},
		}

		if err := rows.Scan(&invoice.ID, &invoice.Server, &invoice.Time, &invoice.Level, &invoice.Hostname, &invoice.ResponseTime, &invoice.Message, &invoice.InvoiceResponse.StatusCode, &invoice.InvoiceRequest.Method, &invoice.InvoiceRequest.URL, &invoice.InvoiceRequest.RemoteAddress, &invoice.InvoiceRequest.RemotePort); err != nil {
			return invoices, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return invoices, err
	}

	return invoices, nil
}
