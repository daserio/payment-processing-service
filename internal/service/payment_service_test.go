package service

import (
	"context"
	"sync"
	"testing"

	"github.com/daserboo/payment-processing-service/internal/domain"
	"github.com/daserboo/payment-processing-service/internal/idempotency"
	"github.com/daserboo/payment-processing-service/internal/locking"
	"github.com/daserboo/payment-processing-service/internal/repository/memory"
)

func TestPaymentService_NoDoubleSpend(t *testing.T) {
	ledger := memory.NewLedgerRepository()
	locker := locking.NewAccountLocker()
	idem := idempotency.NewMemoryStore()

	service := NewPaymentService(ledger, locker, idem)

	// initial balance
	err := ledger.Append(domain.LedgerEntry{
		AccountID: "acc-1",
		Amount:    1000,
		Type:      domain.EntryCredit,
	})
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	errors := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := service.ProcessPayment(
				context.Background(),
				"idem-key-1", // same idempotency key
				"acc-1",
				800,
				"USD",
			)
			errors <- err
		}(i)
	}

	wg.Wait()
	close(errors)

	success := 0
	for err := range errors {
		if err == nil {
			success++
		}
	}

	if success != 1 {
		t.Fatalf("expected exactly 1 successful payment, got %d", success)
	}
}
