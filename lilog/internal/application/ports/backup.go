package ports

import "github.com/frannotsleep/lilog/internal/application/core/domain"

type BackupPort interface {
	Backup(invoices []domain.Invoice) error
}
