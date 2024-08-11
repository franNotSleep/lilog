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

func (app *Application) NewInvoice(server string, invoice domain.Invoice) error {
	err := app.db.Save(server, invoice)
	return err
}

func (app *Application) GetInvoices(server string) ([]domain.Invoice, error) {
	invoices, err := app.db.Get(server)
	return invoices, err
}

func (app *Application) GetServers() ([]string, error) {
	return app.db.Servers()
}
