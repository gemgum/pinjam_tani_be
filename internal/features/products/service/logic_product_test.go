package service_test

import (
	"errors"
	products "projectBE23/internal/features/products"
	"projectBE23/internal/features/products/service"
	"projectBE23/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			UserID:      1,
			Images:      "image.jpg",
			ProductName: "Product1",
			Price:       1000,
			Stock:       10,
			City:        "City1",
			Description: "Description1",
		}

		mockProductData.On("Insert", mockProduct).Return(uint(1), nil).Once()

		productID, err := productService.Create(mockProduct)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), productID)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - empty fields", func(t *testing.T) {
		mockProduct := products.ProductsEntity{}
		productID, err := productService.Create(mockProduct)
		assert.Error(t, err)
		assert.Equal(t, uint(0), productID)
	})

	t.Run("failure - data insert error", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			UserID:      1,
			Images:      "image.jpg",
			ProductName: "Product1",
			Price:       1000,
			Stock:       10,
			City:        "City1",
			Description: "Description1",
		}

		mockProductData.On("Insert", mockProduct).Return(uint(0), assert.AnError).Once()

		productID, err := productService.Create(mockProduct)
		assert.Error(t, err)
		assert.Equal(t, uint(0), productID)
	})
}

func TestGetAllProduct(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success - get all products", func(t *testing.T) {
		mockProducts := []products.ProductsEntity{
			{ProductName: "Product1"},
			{ProductName: "Product2"},
		}

		mockProductData.On("GetAll", 1, 10).Return(mockProducts, nil).Once()

		products, err := productService.GetAllProduct("", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, mockProducts, products)
		mockProductData.AssertExpectations(t)
	})

	t.Run("success - get product by name", func(t *testing.T) {
		mockProducts := []products.ProductsEntity{
			{ProductName: "Product1"},
		}

		mockProductData.On("GetProductByName", "Product1").Return(mockProducts, nil).Once()

		products, err := productService.GetAllProduct("Product1", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, mockProducts, products)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - product by name error", func(t *testing.T) {
		mockProductData.On("GetProductByName", "Product1").Return(nil, assert.AnError).Once()

		products, err := productService.GetAllProduct("Product1", 1, 10)
		assert.Error(t, err)
		assert.Nil(t, products)
	})

	t.Run("failure - GetAll error", func(t *testing.T) {
		mockProductData.On("GetAll", 1, 10).Return(nil, assert.AnError).Once()

		products, err := productService.GetAllProduct("", 1, 10)
		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestGetProductById(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success", func(t *testing.T) {
		mockProduct := &products.ProductsEntity{ProductName: "Product1"}

		mockProductData.On("GetProductById", uint(1)).Return(mockProduct, nil).Once()

		product, err := productService.GetById(1)
		assert.NoError(t, err)
		assert.Equal(t, mockProduct, product)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - invalid ID", func(t *testing.T) {
		product, err := productService.GetById(0)
		assert.Error(t, err)
		assert.Nil(t, product)
	})
}

func TestGetUserProducts(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success", func(t *testing.T) {
		mockProducts := []products.ProductsEntity{
			{ProductName: "Product1"},
			{ProductName: "Product2"},
		}

		mockProductData.On("GetUserProducts", uint(1), 1, 10).Return(mockProducts, nil).Once()

		products, err := productService.GetUserProducts(1, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, mockProducts, products)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - get user products error", func(t *testing.T) {
		mockProductData.On("GetUserProducts", uint(1), 1, 10).Return(nil, assert.AnError).Once()

		products, err := productService.GetUserProducts(1, 1, 10)
		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestUpdateProduct(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			Images:      "hhtps://images",
			ProductName: "UpdatedProduct",
			Price:       20000,
			Stock:       10,
			City:        "City1",
			Description: "Updated Description",
		}

		mockProductData.On("SelectById", uint(1)).Return(&products.ProductsEntity{UserID: 1}, nil).Once()
		mockProductData.On("Update", uint(1), mockProduct).Return(nil).Once()

		err := productService.Update(1, 1, mockProduct)
		assert.NoError(t, err)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - empty fields", func(t *testing.T) {
		emptyFields := []products.ProductsEntity{
			{Images: "", ProductName: "UpdatedProduct", Price: 20000, Stock: 10, City: "City1", Description: "Updated Description"},
			{Images: "https://images", ProductName: "", Price: 20000, Stock: 10, City: "City1", Description: "Updated Description"},
			{Images: "https://images", ProductName: "UpdatedProduct", Price: 0, Stock: 10, City: "City1", Description: "Updated Description"},
			{Images: "https://images", ProductName: "UpdatedProduct", Price: 20000, Stock: 0, City: "City1", Description: "Updated Description"},
			{Images: "https://images", ProductName: "UpdatedProduct", Price: 20000, Stock: 10, City: "", Description: "Updated Description"},
			{Images: "https://images", ProductName: "UpdatedProduct", Price: 20000, Stock: 10, City: "City1", Description: ""},
		}

		for _, mockProduct := range emptyFields {
			err := productService.Update(1, 1, mockProduct)
			assert.Error(t, err)
			assert.Equal(t, "product images/name/price/stock/city/description cannot be empty", err.Error())
		}
	})

	t.Run("failure - user id not match", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			Images:      "https://images",
			ProductName: "UpdatedProduct",
			Price:       20000,
			Stock:       10,
			City:        "City1",
			Description: "Updated Description",
		}

		mockProductData.On("SelectById", uint(1)).Return(&products.ProductsEntity{UserID: 2}, nil).Once()

		err := productService.Update(1, 1, mockProduct)
		assert.Error(t, err)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - select by id error", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			Images:      "https://images",
			ProductName: "UpdatedProduct",
			Price:       20000,
			Stock:       10,
			City:        "City1",
			Description: "Updated Description",
		}

		mockProductData.On("SelectById", uint(1)).Return(nil, errors.New("select error")).Once()

		err := productService.Update(1, 1, mockProduct)
		assert.Error(t, err)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - invalid product id", func(t *testing.T) {
		mockProduct := products.ProductsEntity{
			Images:      "https://images",
			ProductName: "UpdatedProduct",
			Price:       20000,
			Stock:       10,
			City:        "City1",
			Description: "Updated Description",
		}

		err := productService.Update(0, 1, mockProduct)
		assert.Error(t, err)
		assert.Equal(t, "invalid product ID", err.Error())
	})
}

func TestDeleteProduct(t *testing.T) {
	mockProductData := new(mocks.DataProductInterface)
	productService := service.New(mockProductData)

	t.Run("success", func(t *testing.T) {
		mockProductData.On("SelectById", uint(1)).Return(&products.ProductsEntity{UserID: 1}, nil).Once()
		mockProductData.On("Delete", uint(1)).Return(nil).Once()

		err := productService.Delete(1, 1)
		assert.NoError(t, err)
		mockProductData.AssertExpectations(t)
	})

	t.Run("failure - invalid product ID", func(t *testing.T) {
		err := productService.Delete(0, 1)
		assert.Error(t, err)
	})

	t.Run("failure - select by id error", func(t *testing.T) {
		mockProductData.On("SelectById", uint(1)).Return(nil, assert.AnError).Once()

		err := productService.Delete(1, 1)
		assert.Error(t, err)
	})

	t.Run("failure - user id not match", func(t *testing.T) {
		mockProductData.On("SelectById", uint(1)).Return(&products.ProductsEntity{UserID: 2}, nil).Once()

		err := productService.Delete(1, 1)
		assert.Error(t, err)
	})
}
