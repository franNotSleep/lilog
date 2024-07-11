package ports

import "github.com/frannotsleep/lilog/internal/application/core/domain"

type DBPort interface {
	Save(invoice domain.Invoice) error
	Get(pid int32) ([]domain.Invoice, error)
}
