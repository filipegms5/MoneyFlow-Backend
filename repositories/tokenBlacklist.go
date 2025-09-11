package repositories

import (
	"sync"
	"time"
)

// TokenBlacklistRepository stores revoked JTIs in memory with expiration.
type TokenBlacklistRepository struct {
	mu      sync.RWMutex
	entries map[string]time.Time
	quit    chan struct{}
}

// NewTokenBlacklistRepository creates an in-memory blacklist.
// cleanupInterval: how often to sweep expired entries (pass 0 for default 1m).
func NewTokenBlacklistRepository(cleanupInterval time.Duration) *TokenBlacklistRepository {
	if cleanupInterval <= 0 {
		cleanupInterval = time.Minute
	}
	r := &TokenBlacklistRepository{
		entries: make(map[string]time.Time),
		quit:    make(chan struct{}),
	}
	go r.cleanupLoop(cleanupInterval)
	return r
}

func (r *TokenBlacklistRepository) Add(jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	// treat already-expired tokens as short-lived entries
	if ttl <= 0 {
		expiresAt = time.Now().Add(time.Minute)
	}
	r.mu.Lock()
	r.entries[jti] = expiresAt
	r.mu.Unlock()
	return nil
}

func (r *TokenBlacklistRepository) IsRevoked(jti string) (bool, error) {
	r.mu.RLock()
	exp, ok := r.entries[jti]
	r.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		// expired — remove and report not revoked
		r.mu.Lock()
		delete(r.entries, jti)
		r.mu.Unlock()
		return false, nil
	}
	return true, nil
}

// Stop the cleanup goroutine when shutting down (optional).
func (r *TokenBlacklistRepository) Close() {
	close(r.quit)
}

func (r *TokenBlacklistRepository) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			now := time.Now()
			r.mu.Lock()
			for k, exp := range r.entries {
				if now.After(exp) {
					delete(r.entries, k)
				}
			}
			r.mu.Unlock()
		case <-r.quit:
			return
		}
	}
}
