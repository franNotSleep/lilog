package api

import (
	"log"
	"time"
	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Application struct {
	db     ports.DBPort
	backup ports.BackupPort
}

func NewApplication(db ports.DBPort, backup ports.BackupPort) *Application {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:

				invoices, err := db.Export()
				if err != nil {
					log.Printf("Could not export: %v", err)
				}

				err = backup.Backup(invoices)
				if err != nil {
					log.Printf("Could not backup: %v", err)
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

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
