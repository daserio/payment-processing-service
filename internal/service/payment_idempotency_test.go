package service

import (
	"context"
	"testing"

	"github.com/daserboo/payment-processing-service/internal/domain"
	"github.com/daserboo/payment-processing-service/internal/idempotency"
	"github.com/daserboo/payment-processing-service/internal/locking"
	"github.com/daserboo/payment-processing-service/internal/repository/memory"
)

func TestPaymentService_Idempotency(t *testing.T) {
	ledger := memory.NewLedgerRepository()
	locker := locking.NewAccountLocker()
	idem := idempotency.NewMemoryStore()

	svc := NewPaymentService(ledger, locker, idem)

	_ = ledger.Append(domain.LedgerEntry{
		AccountID: "acc-1",
		Amount:    1000,
	})

	key := "idem-key-1"

	err1 := svc.ProcessPayment(context.Background(), key, "acc-1", 500, "USD")
	err2 := svc.ProcessPayment(context.Background(), key, "acc-1", 500, "USD")

	if err1 != nil || err2 != nil {
		t.Fatalf("expected both calls to succeed once, got %v %v", err1, err2)
	}

	entries, _ := ledger.GetByAccount("acc-1")
	if len(entries) != 2 { // initial credit + one debit
		t.Fatalf("expected 2 ledger entries, got %d", len(entries))
	}
}
