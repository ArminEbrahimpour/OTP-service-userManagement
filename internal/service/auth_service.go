package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"otp-service/internal/models"
	"otp-service/internal/repository"
	"otp-service/pkg/jwt"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	SendOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, code string) (*models.AuthResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	otpRepo   repository.OTPRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, otpRepo repository.OTPRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) SendOTP(phoneNumber string) error {
	// Check rate limiting
	count, err := s.otpRepo.IncrementRateLimit(phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %v", err) // Better error formatting
	}

	if count > 3 {
		rateLimit, _ := s.otpRepo.GetRateLimit(phoneNumber)
		if rateLimit != nil {
			// Fixed: Format time properly using RFC3339
			return fmt.Errorf("rate limit exceeded. Try again after %s", rateLimit.ResetTime.Format(time.RFC3339))
		}
		return fmt.Errorf("rate limit exceeded")
	}

	// Generate OTP
	otp, err := s.generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %v", err) // Better error formatting
	}

	// Store OTP with 2-minute expiration
	expiresAt := time.Now().Add(2 * time.Minute)
	if err := s.otpRepo.StoreOTP(phoneNumber, otp, expiresAt); err != nil {
		return fmt.Errorf("failed to store OTP: %v", err) // Better error formatting
	}

	return nil
}

func (s *authService) VerifyOTP(phoneNumber, code string) (*models.AuthResponse, error) {
	// Get and verify OTP
	storedOTP, err := s.otpRepo.GetOTP(phoneNumber)
	if err != nil {
		// Fixed: More user-friendly error message
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	if storedOTP.Code != code {
		return nil, fmt.Errorf("invalid OTP code")
	}

	// Delete used OTP
	s.otpRepo.DeleteOTP(phoneNumber)

	// Check if user exists
	user, err := s.userRepo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Register new user
			user, err = s.userRepo.CreateUser(phoneNumber)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user: %v", err)
		}
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err) // Fixed: Include error details
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) generateOTP() (string, error) {
	// Generate a 6-digit OTP
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return "", err
	}

	n.Add(n, min)
	return n.String(), nil
}
