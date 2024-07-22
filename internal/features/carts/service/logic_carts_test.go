package service_test

import (
	"errors"
	"projectBE23/internal/features/Carts/service"
	carts "projectBE23/internal/features/carts"
	"projectBE23/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCart(t *testing.T) {
	qry := mocks.NewDataCarttInterface(t)
	srv := service.New(qry)
	input := carts.CartEntity{UserID: 1, ProductID: 1, Quantity: 2, TotalPrice: 100000}

	t.Run("Error From Validate", func(t *testing.T) {
		data := carts.CartEntity{UserID: 1, ProductID: 0, Quantity: 0, TotalPrice: 0}
		_, err := srv.Create(data)

		qry.AssertExpectations(t)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cart produts_id/quantity cannot be empty")
	})

	t.Run("Success Create Cart", func(t *testing.T) {
		qry.On("Insert", input).Return(input, nil).Once()
		data, err := srv.Create(input)
		qry.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, input, data)
	})

	t.Run("Error From Query", func(t *testing.T) {
		expectedErr := errors.New("query error")
		qry.On("Insert", input).Return(input, expectedErr).Once()

		data, err := srv.Create(input)

		qry.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, carts.CartEntity{}, data)
	})
}

func TestGetAllCart(t *testing.T) {
	qry := mocks.NewDataCarttInterface(t)
	srv := service.New(qry)
	result := []carts.CartEntity{
		{
			CartID:     1,
			UserID:     1,
			ProductID:  1,
			Quantity:   2,
			TotalPrice: 100000,
		},
	}

	t.Run("Success Get All Cart", func(t *testing.T) {
		qry.On("GetAll").Return(result, nil).Once()
		data, err := srv.GetAllCart()

		qry.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, result, data)
	})
}

func TestFindById(t *testing.T) {
	qry := mocks.NewDataCarttInterface(t)
	srv := service.New(qry)
	result := carts.CartEntity{
		CartID:     1,
		UserID:     1,
		ProductID:  1,
		Quantity:   2,
		TotalPrice: 100000,
	}

	t.Run("Success Get Cart", func(t *testing.T) {
		qry.On("GetById", uint(1)).Return(result, nil).Once()
		data, err := srv.FindById(uint(1))
		qry.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, result, data)
	})
}

func TestUpdateCart(t *testing.T) {
	qry := mocks.NewDataCarttInterface(t)
	srv := service.New(qry)

	t.Run("Success Update Cart", func(t *testing.T) {
		id := uint(1)
		data := carts.CartEntity{UserID: uint(2), ProductID: 1, Quantity: 2, TotalPrice: uint(100000)}

		qry.On("Update", id, data).Return(nil).Once()
		err := srv.Update(id, data)

		qry.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Error From Query", func(t *testing.T) {
		id := uint(1)
		data := carts.CartEntity{UserID: uint(2), ProductID: 1, Quantity: 2, TotalPrice: uint(100000)}

		expectedErr := errors.New("query error")
		qry.On("Update", id, data).Return(expectedErr).Once()

		err := srv.Update(id, data)

		qry.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestDeleteCart(t *testing.T) {
	qry := mocks.NewDataCarttInterface(t)
	srv := service.New(qry)

	t.Run("Success Delete Cart", func(t *testing.T) {
		id := uint(1)
		qry.On("Delete", id).Return(nil).Once()
		err := srv.Delete(id)

		qry.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Error From Query", func(t *testing.T) {
		id := uint(1)
		expectedErr := errors.New("query error")
		qry.On("Delete", id).Return(expectedErr).Once()

		err := srv.Delete(id)

		qry.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
