package ports

import "github.com/frannotsleep/lilog/internal/application/core/domain"

type APIPort interface {
	NewInvoice(domain.Invoice) error
	GetInvoices(pid int32) ([]domain.Invoice, error)
}
