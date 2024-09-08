package api

import (
	"fmt"
	"log"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Application struct {
	db     ports.DBPort
	backup ports.BackupPort
}

func NewApplication(db ports.DBPort, backup ports.BackupPort) *Application {
	app := &Application{
		db:     db,
		backup: backup,
	}

	go func() {
		for {
			select {
      case <-app.backup.C():
				invoices, err := app.db.Export()
				if err != nil {
					log.Println(err)
					continue
				}

				if err := app.backup.Backup(invoices); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

	return app
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

func (app Application) Backup() error {
	invoices, err := app.db.Export()
	if err != nil {
		return fmt.Errorf("Could not export: %v", err)
	}

	err = app.backup.Backup(invoices)
	if err != nil {
		return fmt.Errorf("Could not backup: %v", err)
	}

	return nil
}
