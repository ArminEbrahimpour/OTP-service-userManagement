package service

import (
	"math"
	"otp-service/internal/models"
	"otp-service/internal/repository"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetUsers(page, pageSize int, search string) (*models.UserListResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}
func (s *userService) GetUsers(page, pageSize int, search string) (*models.UserListResponse, error) {
	// default pagination vals
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, total, err := s.userRepo.GetUsers(offset, pageSize, search)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total)) / float64(pageSize))

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
