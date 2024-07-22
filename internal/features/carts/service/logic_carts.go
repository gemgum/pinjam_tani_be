package service

import (
	"errors"
	carts "projectBE23/internal/features/carts"
)

type cartService struct {
	cartData carts.DataCarttInterface
}

func New(cr carts.DataCarttInterface) carts.ServiceCartInterface {
	return &cartService{
		cartData: cr,
	}
}

func (c *cartService) Create(cart carts.CartEntity) (carts.CartEntity, error) {
	if cart.ProductID == 0 || cart.Quantity == 0 {
		return carts.CartEntity{}, errors.New("cart produts_id/quantity cannot be empty")
	}
	data, err := c.cartData.Insert(cart)
	if err != nil {
		return carts.CartEntity{}, err
	}
	return data, nil
}

func (c *cartService) GetAllCart() ([]carts.CartEntity, error) {
	return c.cartData.GetAll()
}

func (c *cartService) FindById(id uint) (carts.CartEntity, error) {
	return c.cartData.GetById(id)
}

func (c *cartService) Update(id uint, cart carts.CartEntity) error {
	err := c.cartData.Update(id, cart)
	if err != nil {
		return err
	}
	return nil
}

func (c *cartService) Delete(id uint) error {
	err := c.cartData.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
