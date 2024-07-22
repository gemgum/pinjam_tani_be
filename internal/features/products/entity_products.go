package products

type ProductsEntity struct {
	ProductID   uint
	UserID      uint
	Images      string
	ProductName string
	Price       int
	Stock       int
	City        string
	Description string
	Quantity    int
}

type DataProductInterface interface {
	Insert(product ProductsEntity) (uint, error)
	Delete(id uint) error
	Update(id uint, product ProductsEntity) error
	GetAll(page, pageSize int) ([]ProductsEntity, error)
	GetProductById(id uint) (*ProductsEntity, error)
	GetProductByName(productName string) ([]ProductsEntity, error)
	GetUserProducts(userID uint, page, pageSize int) ([]ProductsEntity, error)
	SelectById(id uint) (*ProductsEntity, error)
	DecreaseProduct(id uint, sum uint) error
}

type ServiceProductInterface interface {
	Create(product ProductsEntity) (uint, error)
	Delete(id uint, userid uint) error
	Update(id uint, userid uint, product ProductsEntity) error
	GetById(id uint) (product *ProductsEntity, err error)
	GetAllProduct(productName string, page, pageSize int) ([]ProductsEntity, error)
	GetUserProducts(userID uint, page, pageSize int) ([]ProductsEntity, error)
}
