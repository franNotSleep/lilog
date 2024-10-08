package backup

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
)

type BackupAdapter struct {
	interval time.Duration
	out      io.Writer
	ticker   *time.Ticker
}

func NewBackupAdapter(interval time.Duration, out io.Writer) *BackupAdapter {
	return &BackupAdapter{interval: interval, out: out, ticker: time.NewTicker(interval)}
}

func (b BackupAdapter) C() <-chan time.Time {
	return b.ticker.C
}

func (b BackupAdapter) Backup(invoices []domain.Invoice) error {
  b.ticker.Reset(b.interval)

	log.Printf("\033[33mStarting to write backup... 🧾\033[0m\n")
	if len(invoices) == 0 {
		log.Printf("\033[32m[Nothing to Backup] Backup has been successfully written ✅\033[0m\n")
		return nil
	}

	data, err := json.Marshal(invoices)
	if err != nil {
		return err
	}

	_, err = b.out.Write(data)

	if err != nil {
		return err
	}


	log.Printf("\033[32mBackup has been successfully written ✅\033[0m\n")
	return nil
}
