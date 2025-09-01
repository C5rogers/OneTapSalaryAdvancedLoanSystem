package db

import (
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
)

func (db *Database) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) FindUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}
