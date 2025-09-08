package db

import "github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"

func (db *Database) FindCustomerByAccountNumber(acc string) (*models.Customer, error) {
	var c models.Customer
	if err := db.DB.Where("account_number = ?", acc).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *Database) CreateCustomer(c *models.Customer) error {
	return db.DB.Create(c).Error
}

func (db *Database) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	if err := db.DB.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
