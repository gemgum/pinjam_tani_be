package repository

import (
	carts "projectBE23/internal/features/carts"

	"gorm.io/gorm"
)

type cartQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) carts.DataCarttInterface {
	return &cartQuery{
		db: db,
	}
}

func (c *cartQuery) Insert(cart carts.CartEntity) (carts.CartEntity, error) {
	cartGorm := Carts{
		UserID:     cart.UserID,
		ProductID:  cart.ProductID,
		Quantity:   cart.Quantity,
		TotalPrice: cart.TotalPrice,
	}
	tx := c.db.Create(&cartGorm)
	if tx.Error != nil {
		return carts.CartEntity{}, tx.Error
	}
	resCart := carts.CartEntity{
		CartID:     cartGorm.ID,
		UserID:     cartGorm.UserID,
		ProductID:  cartGorm.ProductID,
		Quantity:   cart.Quantity,
		TotalPrice: cart.TotalPrice,
	}
	return resCart, nil
}

func (c *cartQuery) GetAll() ([]carts.CartEntity, error) {
	var allCart []Carts
	tx := c.db.Find(&allCart)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var allCartCore []carts.CartEntity
	for _, v := range allCart {
		allCartCore = append(allCartCore, carts.CartEntity{
			CartID:     v.ID,
			UserID:     v.UserID,
			ProductID:  v.ProductID,
			Quantity:   v.Quantity,
			TotalPrice: v.TotalPrice,
		})
	}

	return allCartCore, nil
}

func (c *cartQuery) GetById(id uint) (carts.CartEntity, error) {
	var cart Carts
	tx := c.db.Where("id = ?", id).First(&cart)
	if tx.Error != nil {
		return carts.CartEntity{}, tx.Error
	}

	cartCore := carts.CartEntity{
		CartID:     cart.ID,
		UserID:     cart.UserID,
		ProductID:  cart.ProductID,
		Quantity:   cart.Quantity,
		TotalPrice: cart.TotalPrice,
	}
	return cartCore, nil
}

func (c *cartQuery) Update(id uint, cart carts.CartEntity) error {
	var cartGorm Carts
	tx := c.db.First(&cartGorm, id)
	if tx.Error != nil {
		return tx.Error
	}
	cartGorm.ProductID = cart.ProductID
	cartGorm.Quantity = cart.Quantity
	cartGorm.TotalPrice = cart.TotalPrice

	tx = c.db.Save(&cartGorm)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// func (c *cartQuery) SelectById(id uint) (*carts.CartEntity, error) {
// 	var cartGorm Carts
// 	tx := c.db.First(&cartGorm, id)
// 	if tx.Error != nil {
// 		return nil, tx.Error
// 	}

// 	var cartcore = carts.CartEntity{
// 		CartID: cartGorm.ID,
// 		UserID: cartGorm.UserID,
// 		// ProducID:   cartGorm.ProductID,
// 		Quantity:   cartGorm.Quantity,
// 		TotalPrice: cartGorm.TotalPrice,
// 	}

// 	return &cartcore, nil
// }

func (c *cartQuery) Delete(id uint) error {
	tx := c.db.Delete(&Carts{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
