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

	// check ratelimit

	count, err := s.otpRepo.IncrementRateLimit(phoneNumber)
	if err != nil {
		return err
	}
	if count > 3 {
		ratelimit, _ := s.otpRepo.GetRateLimit(phoneNumber)
		if ratelimit != nil {
			return fmt.Errorf("ratelimit exceeded try again after %s", ratelimit.ResetTime)
		}
		return fmt.Errorf("ratelimit exceeded")
	}

	// generate new otp
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	//store otp with 2- minute expiration
	expiresAt := time.Now().Add(2 * time.Minute)
	if err := s.otpRepo.StoreOTP(phoneNumber, otp, expiresAt); err != nil {
		return fmt.Errorf("Failed to store OTP:%s", err)
	}
	return nil

}

func (s *authService) VerifyOTP(phoneNumber, code string) (*models.AuthResponse, error) {
	// get and verify otp
	storedotp, err := s.otpRepo.GetOTP(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed getting otp :%s", err)
	}
	if storedotp.Code != code {
		return nil, fmt.Errorf("invalid OTP code")
	}

	// delete used otp
	s.otpRepo.DeleteOTP(phoneNumber)

	// check if user exists
	user, err := s.userRepo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// register new user
			user, err = s.userRepo.CreateUser(phoneNumber)
			if err != nil {
				return nil, fmt.Errorf("Failed creating new user: %v", err)
			}

		} else {
			return nil, fmt.Errorf("Failed to get user: %v", err)
		}
	}
	// generate jwt token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("Failed generating token for user ")
	}
	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) generateOTP() (string, error) {
	// generate 6digit otp
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return "", err
	}
	n.Add(n, min)
	return n.String(), nil
}
