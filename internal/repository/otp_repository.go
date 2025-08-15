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
	CleanupExpiredOTP() error
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
	// printing otp in console as said
	fmt.Printf("OTP for %s : %s (expires at %s) \n", phoneNumber, code, expiresAt)
	return nil
}

func (r *otpRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
	ratelimit, err := r.GetRateLimit(phoneNumber)
	if err != nil {
		// create new rate limit
		r.rateLimits[phoneNumber] = &models.RateLimit{
			PhoneNumber: phoneNumber,
			Count:       1,
			ResetTime:   now.Add(10 * time.Minute),
		}
		return 1, nil
	}
	// increment 1
	ratelimit.Count += 1
	return ratelimit.Count, nil
}

func (r *otpRepository) GetRateLimit(phoneNumber string) (*models.RateLimit, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ratelimit, exists := r.rateLimits[phoneNumber]
	if !exists {
		return nil, fmt.Errorf("rateLimit not exists")
	}
	if time.Now().After(ratelimit.ResetTime) {
		return nil, fmt.Errorf("rate limit expired")
	}
	return ratelimit, nil
}

func (r *otpRepository) CleanupExpiredOTP() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for phoneNumber, otp := range r.otps {
		if now.After(otp.ExpiresAt) {
			delete(r.otps, phoneNumber)
		}
	}
	// cleanup expired ratelimits
	for phoneNumber, ratelimit := range r.rateLimits {
		if now.After(ratelimit.ResetTime) {
			delete(r.rateLimits, phoneNumber)
		}
	}
	return nil
}

func (r *otpRepository) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.CleanupExpiredOTP()
		}
	}
}
