package ports

import "github.com/frannotsleep/lilog/internal/application/core/domain"

type DBPort interface {
	Save(server string, invoice domain.Invoice) error
	Get(server string) ([]domain.Invoice, error)
	Servers() ([]string, error)
}
