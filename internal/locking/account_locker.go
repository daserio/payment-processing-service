package locking

import "sync"

type AccountLocker struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewAccountLocker() *AccountLocker {
	return &AccountLocker{
		locks: make(map[string]*sync.Mutex),
	}
}

func (l *AccountLocker) Lock(accountID string) {
	l.mu.Lock()
	lock, ok := l.locks[accountID]
	if !ok {
		lock = &sync.Mutex{}
		l.locks[accountID] = lock
	}
	l.mu.Unlock()

	lock.Lock()
}

func (l *AccountLocker) Unlock(accountID string) {
	l.mu.Lock()
	lock := l.locks[accountID]
	l.mu.Unlock()

	lock.Unlock()
}
