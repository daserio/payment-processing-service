package repository

import "github.com/daserboo/payment-processing-service/internal/domain"

type LedgerRepository interface {
	Append(entry domain.LedgerEntry) error
	GetByAccount(accountID string) ([]domain.LedgerEntry, error)
}
