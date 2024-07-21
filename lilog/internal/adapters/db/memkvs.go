package db

import "github.com/frannotsleep/lilog/internal/application/core/domain"

type MemKVSAdapter struct {
	kvs map[int32][]domain.Invoice
}

func NewMemKVSAdapter() *MemKVSAdapter {
	return &MemKVSAdapter{
		kvs: make(map[int32][]domain.Invoice),
	}
}

func (m *MemKVSAdapter) Save(invoice domain.Invoice) error {
	if len(m.kvs[invoice.PID]) > 0 {
		m.kvs[invoice.PID] = append(m.kvs[invoice.PID], invoice)
	} else {
		m.kvs[invoice.PID] = []domain.Invoice{invoice}
	}

	return nil
}

func (m *MemKVSAdapter) Get(pid int32) ([]domain.Invoice, error) {
	invoices := m.kvs[pid]
	return invoices, nil
}

func (m *MemKVSAdapter) PIDs() ([]int32, error) {
	keys := make([]int32, 0, len(m.kvs))
	for k := range m.kvs {
		keys = append(keys, k)
	}

	pids := make([]int32, 0, len(m.kvs))
	for pid := range m.kvs {
		pids = append(pids, pid)
	}

	return pids, nil
}
