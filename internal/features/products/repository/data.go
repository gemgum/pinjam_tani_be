package repository

import (
	datacarts "projectBE23/internal/features/carts/repository"
	dataorder "projectBE23/internal/features/orders/repository"

	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	UserID      uint
	Images      string
	ProductName string
	Price       int
	Stock       int
	City        string
	Description string
	Orders      []dataorder.OrderItems `gorm:"foreignKey:ProductID"`
	Carts       []datacarts.Carts      `gorm:"foreignKey:ProductID"`
}
