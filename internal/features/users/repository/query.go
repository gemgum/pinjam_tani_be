package repository

import (
	"log"
	"pinjamtani_project/internal/features/users"

	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) users.QueryUserInterface {
	return &userQuery{
		db: db,
	}
}

// CreateAccount implements users.DataUserInterface.
func (u *userQuery) CreateAccount(account users.User) (uint, error) {
	userGorm := Users{
		Images:      account.Images,
		UserName:    account.UserName,
		Email:       account.Email,
		Password:    account.Password,
		PhoneNumber: account.PhoneNumber,
		Address:     account.Address,
	}
	tx := u.db.Create(&userGorm)

	if tx.Error != nil {
		log.Printf("CreateAccount: Error creating account: %v", tx.Error)
		return 0, tx.Error
	}

	return userGorm.ID, nil
}

// AccountByEmail implements users.DataUserInterface.
func (u *userQuery) AccountByEmail(email string) (*users.User, error) {
	var userData Users
	tx := u.db.Where("email = ?", email).First(&userData)
	if tx.Error != nil {
		log.Printf("AccountByEmail: Error finding account: %v", tx.Error)
		return nil, tx.Error
	}
	// mapping
	var users = users.User{
		UserID:      userData.ID,
		UserName:    userData.UserName,
		Email:       userData.Email,
		Password:    userData.Password,
		Address:     userData.Address,
		PhoneNumber: userData.PhoneNumber,
	}

	return &users, nil
}

func (u *userQuery) AccountById(userid uint) (*users.User, error) {
	var userData Users
	tx := u.db.First(&userData, userid)
	if tx.Error != nil {
		log.Printf("AccountById: Error finding account: %v", tx.Error)
		return nil, tx.Error
	}
	// mapping
	var user = users.User{
		UserID:      userData.ID,
		Images:      userData.Images,
		UserName:    userData.UserName,
		Email:       userData.Email,
		PhoneNumber: userData.PhoneNumber,
		Address:     userData.Address,
	}

	return &user, nil
}

func (u *userQuery) UpdateAccount(userid uint, account users.User) error {
	var userGorm Users
	tx := u.db.First(&userGorm, userid)
	if tx.Error != nil {
		log.Printf("UpdateAccount: Error finding account: %v", tx.Error)
		return tx.Error
	}
	userGorm.Images = account.Images
	userGorm.UserName = account.UserName
	userGorm.Email = account.Email
	userGorm.Password = account.Password
	userGorm.PhoneNumber = account.PhoneNumber
	userGorm.Address = account.Address

	tx = u.db.Save(&userGorm)
	if tx.Error != nil {
		log.Printf("UpdateAccount: Error updating account: %v", tx.Error)
		return tx.Error
	}
	return nil
}

func (u *userQuery) DeleteAccount(userid uint) error {
	tx := u.db.Delete(&Users{}, userid)
	if tx.Error != nil {
		log.Printf("DeleteAccount: Error deleting account: %v", tx.Error)
		return tx.Error
	}
	return nil
}
