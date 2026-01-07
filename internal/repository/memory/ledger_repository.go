package memory

import (
	"sync"

	"github.com/daserboo/payment-processing-service/internal/domain"
)

type LedgerRepository struct {
	mu      sync.RWMutex
	entries map[string][]domain.LedgerEntry // accountID -> entries
}

func NewLedgerRepository() *LedgerRepository {
	return &LedgerRepository{
		entries: make(map[string][]domain.LedgerEntry),
	}
}

func (r *LedgerRepository) Append(entry domain.LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries[entry.AccountID] = append(r.entries[entry.AccountID], entry)
	return nil
}

func (r *LedgerRepository) GetByAccount(accountID string) ([]domain.LedgerEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := r.entries[accountID]

	// copy slice to avoid external mutation
	result := make([]domain.LedgerEntry, len(entries))
	copy(result, entries)

	return result, nil
}
