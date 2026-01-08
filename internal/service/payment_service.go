package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/daserboo/payment-processing-service/internal/domain"
	"github.com/daserboo/payment-processing-service/internal/idempotency"
	"github.com/daserboo/payment-processing-service/internal/locking"
	"github.com/daserboo/payment-processing-service/internal/repository"
)

type PaymentService struct {
	ledger repository.LedgerRepository
	locker *locking.AccountLocker
	idem   idempotency.Store
}

func NewPaymentService(
	ledger repository.LedgerRepository,
	locker *locking.AccountLocker,
	idem idempotency.Store,
) *PaymentService {
	return &PaymentService{
		ledger: ledger,
		locker: locker,
		idem:   idem,
	}
}

func (s *PaymentService) ProcessPayment(
	ctx context.Context,
	idempotencyKey string,
	accountID string,
	amount int64,
	currency string,
) error {
	if idempotencyKey != "" {
		if res, ok := s.idem.Get(idempotencyKey); ok {
			return res.Err
		}
	}

	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	s.locker.Lock(accountID)
	defer s.locker.Unlock(accountID)

	entries, err := s.ledger.GetByAccount(accountID)
	if err != nil {
		return err
	}

	balance := domain.CalculateBalance(entries)
	if balance < amount {
		err = domain.ErrInsufficientFunds
		if idempotencyKey != "" {
			s.idem.Set(idempotencyKey, idempotency.Result{
				Err:       err,
				CreatedAt: time.Now(),
			}, 5*time.Minute)
		}
		return err
	}

	entry := domain.LedgerEntry{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Amount:    -amount,
		Currency:  currency,
		Type:      domain.EntryDebit,
		CreatedAt: time.Now(),
	}

	err = s.ledger.Append(entry)

	if idempotencyKey != "" {
		s.idem.Set(idempotencyKey, idempotency.Result{
			TransactionID: entry.ID,
			Err:           err,
			CreatedAt:     time.Now(),
		}, 5*time.Minute)
	}

	return err
}
