package repository

import (
	dataorder "projectBE23/internal/features/orders/repository"

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
