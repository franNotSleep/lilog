package api

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort, backupInterval time.Duration) *Application {
	ticker := time.NewTicker(backupInterval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := Export("logs.txt", db)
				if err != nil {
					log.Printf("Could not write to logs.txt: %v", err)
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

func Export(filename string, db ports.DBPort) error {
  log.Printf("Writing logs to %s...", filename)
	invoices, err := db.Export()
	if err != nil {
		return err
	}

	if len(invoices) == 0 {
		return nil
	}

	b, err := json.Marshal(invoices)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, b, 0666)
	return err
}
