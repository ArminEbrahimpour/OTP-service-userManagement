package models

import "time"

// represent user
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	PhoneNumber      string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	RegistrationDate time.Time `json:"registration_data" gorm:"not null"`
	CreateAt         time.Time `json:"created_at" `
	UpdatedAt        time.Time `json:"updated_at"`
}

// represent otp token
type OTP struct {
	PhoneNumber string    `json:"phone_number" `
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// represent request for otp generation
type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// represent request for otp verification
type OTPVerifyReq struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

// represetn response after successfull authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// represent user list
type UserListResponse struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_page"`
}

// represent rate limit info

type RateLimit struct {
	PhoneNumber string    `json:"phone_number"`
	Count       int       `json:"count"`
	ResetTime   time.Time `json:"reset_time"`
}
