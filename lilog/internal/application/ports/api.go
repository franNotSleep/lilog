package ports

import (
	"time"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
)

type APIPort interface {
	NewInvoice(server string, invoice domain.Invoice) error
	GetInvoices(server string) ([]domain.Invoice, error)
	GetServers() ([]string, error)
	GetBetween(from time.Time, until time.Time) ([]domain.Invoice, error)
	Backup() error
}
