package repository

import (
	"sync"
	"time"

	"github.com/adarsh-jaiss/assingment/api/models"
)

type UserRateLimit struct {
	WindowStart   time.Time
	AcceptedCount int
	RejectedCount int
}

type RateLimitRepo interface {
	CheckAndRecord(userID string) (bool, error)
	GetStats() []models.UserStats
}

type rateLimitRepoImpl struct {
	mu    sync.RWMutex
	users map[string]*UserRateLimit
}

func NewRateLimitRepo() RateLimitRepo {
	return &rateLimitRepoImpl{
		users: make(map[string]*UserRateLimit),
	}
}

func (r *rateLimitRepoImpl) CheckAndRecord(userID string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	user, exists := r.users[userID]
	if !exists {
		user = &UserRateLimit{
			WindowStart:   now,
			AcceptedCount: 0,
			RejectedCount: 0,
		}
		r.users[userID] = user
	}

	// Check if the current 1-minute window has expired
	if now.Sub(user.WindowStart) >= time.Minute {
		user.WindowStart = now
		user.AcceptedCount = 0
	}

	// Check rate limit
	if user.AcceptedCount < 5 {
		user.AcceptedCount++
		return true, nil
	}

	user.RejectedCount++
	return false, nil
}

func (r *rateLimitRepoImpl) GetStats() []models.UserStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make([]models.UserStats, 0, len(r.users))
	for id, user := range r.users {
		stats = append(stats, models.UserStats{
			UserID:        id,
			AcceptedCount: user.AcceptedCount,
			RejectedCount: user.RejectedCount,
		})
	}

	return stats
}
