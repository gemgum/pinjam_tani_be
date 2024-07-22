package repository

import (
	datacarts "projectBE23/internal/features/carts/repository"
	dataorder "projectBE23/internal/features/orders/repository"
	dataproducts "projectBE23/internal/features/products/repository"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Images      string
	UserName    string
	Email       string `gorm:"unique"`
	Password    string
	PhoneNumber string
	Address     string
	Orders      []dataorder.OrderStatus `gorm:"foreignKey:UserID"`
	Products    []dataproducts.Products `gorm:"foreignKey:UserID"`
	Carts       []datacarts.Carts       `gorm:"foreignKey:UserID"`
	OrderItems  []dataorder.OrderItems  `gorm:"foreignKey:UserID"`
}
