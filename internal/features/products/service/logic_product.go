package service

import (
	"errors"
	"log"

	products "projectBE23/internal/features/products"
)

type productService struct {
	productData products.DataProductInterface
}

func New(pr products.DataProductInterface) products.ServiceProductInterface {
	return &productService{
		productData: pr,
	}
}

// Create implements products.ServiceProductsInterface.
func (p *productService) Create(product products.ProductsEntity) (uint, error) {
	if product.Images == "" || product.ProductName == "" || product.Price == 0 || product.Stock == 0 || product.City == "" || product.Description == "" {
		err := errors.New("product images/name/price/stock/city/description cannot be empty")
		log.Println(err)
		return 0, err
	}
	productID, err := p.productData.Insert(product)
	if err != nil {
		log.Println("Error creating product:", err)
		return 0, err
	}
	return productID, nil
}

// Delete implements products.ServiceProductsInterface.
func (p *productService) Delete(id uint, userid uint) error {
	if id <= 0 {
		err := errors.New("invalid products ID")
		log.Println(err)
		return err
	}
	cekuserid, err := p.productData.SelectById(id)
	if err != nil {
		log.Println("Error deleting product:", err)
		return err
	}

	if cekuserid.UserID != userid {
		err := errors.New("user id not match, cannot delete products")
		log.Println(err)
		return err
	}

	return p.productData.Delete(id)
}

// GetAllProduct implements products.ServiceProductsInterface.
func (p *productService) GetAllProduct(productName string, page, pageSize int) ([]products.ProductsEntity, error) {
	if productName != "" {
		// Jika nama product diberikan, lakukan pencarian berdasarkan nama product
		product, err := p.productData.GetProductByName(productName)
		if err != nil {
			log.Println("Error retrieving product by name:", err)
			return nil, err
		}
		return product, nil
	}
	// Jika tidak ada nama product yang diberikan, kembalikan semua product
	allProduct, err := p.productData.GetAll(page, pageSize)
	if err != nil {
		log.Println("Error retrieving all products:", err)
		return nil, err
	}
	return allProduct, nil
}

// GetById implements products.ServiceProductsInterface.
func (p *productService) GetById(id uint) (product *products.ProductsEntity, err error) {
	if id <= 0 {
		err = errors.New("id not valid")
		log.Println(err)
		return nil, err
	}
	return p.productData.GetProductById(id)
}

// Update implements products.ServiceProductsInterface.
func (p *productService) Update(id uint, userid uint, product products.ProductsEntity) error {
	if id == 0 {
		err := errors.New("invalid product ID")
		log.Println(err)
		return err
	}
	if product.Images == "" || product.ProductName == "" || product.Price == 0 || product.Stock == 0 || product.City == "" || product.Description == "" {
		err := errors.New("product images/name/price/stock/city/description cannot be empty")
		log.Println(err)
		return err
	}

	cekuserid, err := p.productData.SelectById(id)
	if err != nil {
		log.Println("Error updating product:", err)
		return err
	}

	if cekuserid.UserID != userid {
		err := errors.New("user id not match, cannot update products")
		log.Println(err)
		return err
	}

	return p.productData.Update(id, product)
}

func (p *productService) GetUserProducts(userID uint, page, pageSize int) ([]products.ProductsEntity, error) {
	userProducts, err := p.productData.GetUserProducts(userID, page, pageSize)
	if err != nil {
		log.Println("Error retrieving user products:", err)
		return nil, err
	}
	return userProducts, nil
}
