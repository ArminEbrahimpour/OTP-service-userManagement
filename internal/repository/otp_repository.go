package repository

import (
	"fmt"
	"otp-service/internal/models"
	"sync"
	"time"
)

type OTPRepository interface {
	StoreOTP(phoneNumber, code string, expiresAt time.Time) error
	GetOTP(phoneNumber string) (*models.OTP, error)
	DeleteOTP(phoneNumber string) error
	IncrementRateLimit(phoneNumber string) (int, error)
	GetRateLimit(phoneNumber string) (*models.RateLimit, error)
	CleanupExpiredOTPs() // Fixed: Method name should match implementation
}

type otpRepository struct {
	otps       map[string]*models.OTP
	rateLimits map[string]*models.RateLimit
	mu         sync.RWMutex
}

func NewOTPRepository() OTPRepository {
	repo := &otpRepository{
		otps:       make(map[string]*models.OTP),
		rateLimits: make(map[string]*models.RateLimit),
	}

	// Start cleanup routine
	go repo.cleanupRoutine()

	return repo
}

func (r *otpRepository) StoreOTP(phoneNumber, code string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.otps[phoneNumber] = &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}

	// Print OTP to console as required
	fmt.Printf("OTP for %s: %s (expires at %s)\n", phoneNumber, code, expiresAt.Format(time.RFC3339))

	return nil
}

func (r *otpRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	r.mu.RLock() // Fixed: Use RLock for read operations
	defer r.mu.RUnlock()

	otp, exists := r.otps[phoneNumber]
	if !exists {
		return nil, fmt.Errorf("OTP not found")
	}

	if time.Now().After(otp.ExpiresAt) {
		return nil, fmt.Errorf("OTP expired")
	}

	return otp, nil
}

func (r *otpRepository) DeleteOTP(phoneNumber string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.otps, phoneNumber)
	return nil
}

func (r *otpRepository) IncrementRateLimit(phoneNumber string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	rateLimit, exists := r.rateLimits[phoneNumber]

	if !exists || now.After(rateLimit.ResetTime) {
		// Create new rate limit entry
		r.rateLimits[phoneNumber] = &models.RateLimit{
			PhoneNumber: phoneNumber,
			Count:       1,
			ResetTime:   now.Add(10 * time.Minute),
		}
		return 1, nil
	}

	// Increment existing count
	rateLimit.Count++
	return rateLimit.Count, nil
}

func (r *otpRepository) GetRateLimit(phoneNumber string) (*models.RateLimit, error) {
	r.mu.RLock() // Fixed: Use RLock for read operations
	defer r.mu.RUnlock()

	rateLimit, exists := r.rateLimits[phoneNumber]
	if !exists {
		return nil, fmt.Errorf("rate limit not found")
	}

	if time.Now().After(rateLimit.ResetTime) {
		return nil, fmt.Errorf("rate limit expired")
	}

	return rateLimit, nil
}

func (r *otpRepository) CleanupExpiredOTPs() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Clean up expired OTPs
	for phoneNumber, otp := range r.otps {
		if now.After(otp.ExpiresAt) {
			delete(r.otps, phoneNumber)
		}
	}

	// Clean up expired rate limits
	for phoneNumber, rateLimit := range r.rateLimits {
		if now.After(rateLimit.ResetTime) {
			delete(r.rateLimits, phoneNumber)
		}
	}
}

func (r *otpRepository) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.CleanupExpiredOTPs()
		}
	}
}
