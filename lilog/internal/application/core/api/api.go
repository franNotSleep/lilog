package api

import (
	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (app *Application) NewInvoice(invoice domain.Invoice) error {
	err := app.db.Save(invoice)
	return err
}

func (app *Application) GetInvoices(pid int64) ([]domain.Invoice, error) {
	invoices, err := app.db.Get(pid)
	return invoices, err
}
