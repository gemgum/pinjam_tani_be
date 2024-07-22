package repository

import (
	datacarts "pinjamtani_project/internal/features/carts/repository"
	dataorder "pinjamtani_project/internal/features/orders/repository"

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
