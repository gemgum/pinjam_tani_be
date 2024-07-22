package carts

type CartEntity struct {
	CartID     uint
	UserID     uint
	ProductID  uint
	Quantity   uint
	TotalPrice uint
}

type QueryCartInterface interface {
	Insert(cart CartEntity) (CartEntity, error)
	Delete(id uint) error
	Update(id uint, cart CartEntity) error
	GetAll() ([]CartEntity, error)
	GetById(id uint) (CartEntity, error)
}

type ServiceCartInterface interface {
	Create(cart CartEntity) (CartEntity, error)
	Delete(id uint) error
	Update(id uint, cart CartEntity) error
	GetAllCart() ([]CartEntity, error)
	FindById(id uint) (CartEntity, error)
}
