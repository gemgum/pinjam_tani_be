package repository

import (
	dataorder "pinjamtani_project/internal/features/orders/repository"

	"gorm.io/gorm"
)

type Carts struct {
	gorm.Model
	UserID     uint
	ProductID  uint
	Quantity   uint
	TotalPrice uint
	Orders     []dataorder.OrderItems `gorm:"foreignKey:CartID"`
}
