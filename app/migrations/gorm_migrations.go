package migrations

import (
	datacarts "projectBE23/internal/features/carts/repository"
	dataorder "projectBE23/internal/features/orders/repository"
	dataproducts "projectBE23/internal/features/products/repository"
	datausers "projectBE23/internal/features/users/repository"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&datausers.Users{})
	db.AutoMigrate(&dataorder.OrderItems{})
	db.AutoMigrate(&dataorder.OrderStatus{})
	db.AutoMigrate(&datacarts.Carts{})
	db.AutoMigrate(&dataproducts.Products{})
}
