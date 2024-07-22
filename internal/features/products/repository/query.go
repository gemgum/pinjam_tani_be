package repository

import (
	"errors"
	"log"
	products "pinjamtani_project/internal/features/products"
	"pinjamtani_project/internal/utils"

	"gorm.io/gorm"
)

type productQuery struct {
	db *gorm.DB
}

func NewPaymentQuery(db *gorm.DB) products.QueryProductInterface {
	return &productQuery{
		db: db,
	}
}

// Delete implements products.DataProductInterface.
func (p *productQuery) Delete(id uint) error {
	tx := p.db.Delete(&Products{}, id)
	if tx.Error != nil {
		log.Printf("Error deleting product with id %d: %v", id, tx.Error)
		return tx.Error
	}
	return nil
}

// GetAll implements products.DataProductInterface.
func (p *productQuery) GetAll(page, pageSize int) ([]products.ProductsEntity, error) {
	var allProduct []Products
	pagination := utils.NewPagination(page, pageSize)

	tx := p.db.Limit(pagination.PageSize).Offset(pagination.Offset()).Find(&allProduct)
	if tx.Error != nil {
		log.Printf("Error fetching all products: %v", tx.Error)
		return nil, tx.Error
	}

	var allProductCore []products.ProductsEntity
	for _, v := range allProduct {
		allProductCore = append(allProductCore, products.ProductsEntity{
			ProductID:   v.ID,
			Images:      v.Images,
			ProductName: v.ProductName,
			Price:       v.Price,
			Stock:       v.Stock,
			City:        v.City,
		})
	}

	return allProductCore, nil
}

// GetProductByName implements products.DataProductInterface.
func (p *productQuery) GetProductByName(productName string) ([]products.ProductsEntity, error) {
	var allProduct []Products
	tx := p.db.Where("product_name =?", productName).Find(&allProduct)
	if tx.Error != nil {
		log.Printf("Error fetching product by name %s: %v", productName, tx.Error)
		return nil, tx.Error
	}

	// mapping
	var allProductCore []products.ProductsEntity
	for _, v := range allProduct {
		allProductCore = append(allProductCore, products.ProductsEntity{
			Images:      v.Images,
			ProductName: v.ProductName,
			Price:       v.Price,
			Stock:       v.Stock,
			City:        v.City,
		})
	}

	return allProductCore, nil
}

// Insert implements products.DataProductInterface.
func (p *productQuery) Insert(product products.ProductsEntity) (uint, error) {
	productGorm := Products{
		UserID:      product.UserID,
		Images:      product.Images,
		ProductName: product.ProductName,
		Price:       product.Price,
		Stock:       product.Stock,
		City:        product.City,
		Description: product.Description,
	}
	tx := p.db.Create(&productGorm)

	if tx.Error != nil {
		log.Printf("Error inserting product: %v", tx.Error)
		return 0, tx.Error
	}
	return productGorm.ID, nil
}

// SelectById implements products.DataProductInterface.
func (p *productQuery) SelectById(id uint) (*products.ProductsEntity, error) {
	var productGorm Products
	tx := p.db.First(&productGorm, id)
	if tx.Error != nil {
		log.Printf("Error fetching product by id %d: %v", id, tx.Error)
		return nil, tx.Error
	}

	// mapping
	var productcore = products.ProductsEntity{
		ProductID:   productGorm.ID,
		UserID:      productGorm.UserID,
		ProductName: productGorm.ProductName,
		Price:       productGorm.Price,
		Stock:       productGorm.Stock,
		City:        productGorm.City,
		Description: productGorm.Description,
	}

	return &productcore, nil
}

// Update implements products.DataProductInterface.
func (p *productQuery) Update(id uint, product products.ProductsEntity) error {
	var productGorm Products
	tx := p.db.First(&productGorm, id)
	if tx.Error != nil {
		log.Printf("Error fetching product for update with id %d: %v", id, tx.Error)
		return tx.Error
	}
	productGorm.Images = product.Images
	productGorm.ProductName = product.ProductName
	productGorm.Price = product.Price
	productGorm.Stock = product.Stock
	productGorm.City = product.City
	productGorm.Description = product.Description

	tx = p.db.Save(&productGorm)
	if tx.Error != nil {
		log.Printf("Error updating product with id %d: %v", id, tx.Error)
		return tx.Error
	}
	return nil
}

// GetProductById implements products.DataProductInterface.
func (p *productQuery) GetProductById(id uint) (*products.ProductsEntity, error) {
	var productGorm Products
	tx := p.db.First(&productGorm, id)
	if tx.Error != nil {
		log.Printf("Error fetching product by id %d: %v", id, tx.Error)
		return nil, tx.Error
	}

	// mapping
	var productcore = products.ProductsEntity{
		ProductID:   productGorm.ID,
		UserID:      productGorm.UserID,
		Images:      productGorm.Images,
		ProductName: productGorm.ProductName,
		Price:       productGorm.Price,
		Stock:       productGorm.Stock,
		City:        productGorm.City,
		Description: productGorm.Description,
	}

	return &productcore, nil
}

func (p *productQuery) GetUserProducts(userID uint, page, pageSize int) ([]products.ProductsEntity, error) {
	var userProducts []Products
	pagination := utils.NewPagination(page, pageSize)

	tx := p.db.Where("user_id = ?", userID).Limit(pagination.PageSize).Offset(pagination.Offset()).Find(&userProducts)
	if tx.Error != nil {
		log.Printf("Error fetching products for user id %d: %v", userID, tx.Error)
		return nil, tx.Error
	}

	var userProductsCore []products.ProductsEntity
	for _, v := range userProducts {
		userProductsCore = append(userProductsCore, products.ProductsEntity{
			ProductID:   v.ID,
			UserID:      v.UserID,
			Images:      v.Images,
			ProductName: v.ProductName,
			Price:       v.Price,
			Stock:       v.Stock,
			City:        v.City,
			Description: v.Description,
		})
	}

	return userProductsCore, nil
}

func (p *productQuery) DecreaseProduct(prodID uint, sum uint) error {
	dataProducts, err := p.GetProductById(prodID)
	if err != nil {
		log.Println("error get cart data", err.Error())
		return err
	}
	if uint(dataProducts.Stock) >= sum {
		dataQuantity := uint(dataProducts.Stock) - sum
		tx := p.db.Model(&Products{}).Where("id = ?", prodID).Update("stock", dataQuantity)
		if tx.Error != nil {
			log.Println("error get cart", tx.Error.Error())
			return err
		}

		rowsAffected := tx.RowsAffected
		if rowsAffected <= 0 {
			log.Println("no rows affected")
			return errors.New("no rows affected")

		}
	} else {
		log.Println("error quantity lower than sum")
		return errors.New("error quantity lower than sum")
	}
	return nil
}
