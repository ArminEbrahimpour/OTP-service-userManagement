package repository

import (
	"otp-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(phoneNumber string) (*models.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUsers(offset, limit int, search string) ([]models.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// create new user repo
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// create new user
func (r *userRepository) CreateUser(phoneNumber string) (*models.User, error) {
	user := &models.User{
		PhoneNumber:      phoneNumber,
		RegistrationDate: time.Now(),
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers(offset, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if search == "" {
		query = query.Where("phone_number ILIKE ?", "%"+search+"%s")
	}

	//get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// get users
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
