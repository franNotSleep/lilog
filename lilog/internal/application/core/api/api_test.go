package api

import (
	"testing"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
}

func (db *mockDB) Save(invoice domain.Invoice) error {
	args := db.Called(invoice)
	return args.Error(0)
}

func (db *mockDB) Get(pid int32) ([]domain.Invoice, error) {
	args := db.Called(pid)
	return args.Get(0).([]domain.Invoice), args.Error(1)
}

func TestNewInvoice(t *testing.T) {
	db := new(mockDB)
	db.On("Save", mock.Anything).Return(nil)

	application := NewApplication(db)
	err := application.NewInvoice(domain.Invoice{})

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvoices(t *testing.T) {
	db := new(mockDB)
	inv1 := domain.Invoice{PID: 123}
	inv2 := domain.Invoice{PID: 125}

	db.On("Get", mock.Anything).Return([]domain.Invoice{inv1, inv2}, nil)

	application := NewApplication(db)
	invoices, err := application.db.Get(3122)

	if err != nil {
		t.Error(err)
	}

	if len(invoices) != 2 {
		t.Errorf("Expected invoices length to be %d; expected %d", 2, len(invoices))
	}
}
