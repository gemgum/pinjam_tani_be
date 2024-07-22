package migrations

import (
	datacarts "pinjamtani_project/internal/features/carts/repository"
	dataorder "pinjamtani_project/internal/features/orders/repository"
	dataproducts "pinjamtani_project/internal/features/products/repository"
	datausers "pinjamtani_project/internal/features/users/repository"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&datausers.Users{})
	db.AutoMigrate(&dataorder.OrderItems{})
	db.AutoMigrate(&dataorder.OrderStatus{})
	db.AutoMigrate(&datacarts.Carts{})
	db.AutoMigrate(&dataproducts.Products{})
}
