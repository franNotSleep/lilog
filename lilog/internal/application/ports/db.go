package ports

import (
	"time"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
)

type DBPort interface {
	Save(server string, invoice domain.Invoice) error
	Get(server string) ([]domain.Invoice, error)
	Servers() ([]string, error)
	GetBetween(server string, from time.Time, until time.Time) ([]domain.Invoice, error)
	Export() ([]domain.Invoice, error)
}
