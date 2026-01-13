package memory

import (
	"sync"
	"testing"

	"github.com/daserio/payment-processing-service/internal/domain"
)

func TestLedgerRepository_AppendAndGet(t *testing.T) {
	repo := NewLedgerRepository()

	entry := domain.LedgerEntry{
		AccountID: "acc-1",
		Amount:    1000,
	}

	if err := repo.Append(entry); err != nil {
		t.Fatalf("append failed: %v", err)
	}

	entries, _ := repo.GetByAccount("acc-1")

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestLedgerRepository_ConcurrentAccess(t *testing.T) {
	repo := NewLedgerRepository()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = repo.Append(domain.LedgerEntry{
				AccountID: "acc-1",
				Amount:    10,
			})
		}()
	}

	wg.Wait()

	entries, _ := repo.GetByAccount("acc-1")
	if len(entries) != 100 {
		t.Fatalf("expected 100 entries, got %d", len(entries))
	}
}
