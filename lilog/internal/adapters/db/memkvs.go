package db

import (
	"github.com/frannotsleep/lilog/internal/application/core/domain"
)

type MemKVSAdapter struct {
	kvs map[string][]domain.Invoice
}

func NewMemKVSAdapter() *MemKVSAdapter {
	return &MemKVSAdapter{
		kvs: make(map[string][]domain.Invoice),
	}
}

func (m *MemKVSAdapter) Save(server string, invoice domain.Invoice) error {
	if len(m.kvs[server]) > 0 {
		m.kvs[server] = append(m.kvs[server], invoice)
	} else {
		m.kvs[server] = []domain.Invoice{invoice}
	}

	return nil
}

func (m *MemKVSAdapter) Get(server string) ([]domain.Invoice, error) {
	invoices := m.kvs[server]
	return invoices, nil
}

func (m *MemKVSAdapter) Servers() ([]string, error) {
	keys := make([]string, 0, len(m.kvs))
	for k := range m.kvs {
		keys = append(keys, k)
	}

	servers := make([]string, 0, len(m.kvs))
	for server := range m.kvs {
		servers = append(servers, server)
	}

	return servers, nil
}
