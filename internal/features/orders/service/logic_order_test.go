package service_test

import (
	"errors"
	order "projectBE23/internal/features/orders"
	"projectBE23/internal/features/orders/service"
	"projectBE23/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddOrderItems(t *testing.T) {
	dataOrder := mocks.NewQueryOrderInterface(t)
	midl := mocks.NewMiddlewaresInterface(t)
	ordSrv := service.NewOrderService(dataOrder, midl)

	paymentUrl := uuid.NewString()
	trxDatetime := time.Now().Local()
	inputOrderItem := order.OrderItemEntity{ProductID: 1, CartID: 2, Quantity: 10, UserID: 1}
	inputOrderStatus := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime}

	t.Run("Success Add OrderItem", func(t *testing.T) {
		inputOrderItemQry := order.OrderItemEntity{ProductID: 1, CartID: 2, Quantity: 10, UserID: 1}
		inputOrderStatusQry := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime.Local()}
		dataOrder.On("AddOrderItems", inputOrderItem, inputOrderStatus).Return(inputOrderItemQry.ProductID, inputOrderItemQry.Quantity, nil).Once()
		productID, quantity, err := ordSrv.AddOrderItems(inputOrderItemQry, inputOrderStatusQry)

		assert.Nil(t, err)
		assert.Equal(t, inputOrderItemQry.ProductID, productID)
		assert.Equal(t, inputOrderItemQry.Quantity, quantity)
	})

	t.Run("Failed Add OrderItem", func(t *testing.T) {
		inputOrderItemQry := order.OrderItemEntity{ProductID: 0, CartID: 2, Quantity: 0, UserID: 1}
		inputOrderStatusQry := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime.Local()}
		inputOrderItem.ProductID = 0
		inputOrderItem.Quantity = 0
		dataOrder.On("AddOrderItems", inputOrderItem, inputOrderStatus).Return(inputOrderItemQry.ProductID, inputOrderItemQry.Quantity, errors.New("validasi tidak sesuai")).Once()
		productID, quantity, err := ordSrv.AddOrderItems(inputOrderItemQry, inputOrderStatusQry)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "validasi tidak sesuai")
		assert.Equal(t, inputOrderItemQry.ProductID, productID)
		assert.Equal(t, inputOrderItemQry.Quantity, quantity)
	})
}

func TestGetProductById(t *testing.T) {
	dataOrder := mocks.NewQueryOrderInterface(t)
	midl := mocks.NewMiddlewaresInterface(t)
	ordSrv := service.NewOrderService(dataOrder, midl)

	userID := 1
	dataListOrder := order.ListOrderItemEntity{}
	t.Run("Success Add OrderItem", func(t *testing.T) {
		dataOrder.On("GetOrders", uint(userID)).Return(dataListOrder, nil).Once()
		listOrder, err := ordSrv.GetOrders(uint(userID))

		assert.Nil(t, err)
		assert.Equal(t, listOrder, dataListOrder)
	})

	t.Run("Failed GetOrders", func(t *testing.T) {
		expectedErr := errors.New("data id tidak ditemukan / tidak valid")

		dataOrder.On("GetOrders", uint(userID)).Return(order.ListOrderItemEntity{}, expectedErr).Once()
		listOrder, err := ordSrv.GetOrders(uint(userID))

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, listOrder, order.ListOrderItemEntity{})
	})
}

func TestAddOrderStatuses(t *testing.T) {
	dataOrder := mocks.NewQueryOrderInterface(t)
	midl := mocks.NewMiddlewaresInterface(t)
	ordSrv := service.NewOrderService(dataOrder, midl)

	paymentUrl := uuid.NewString()
	trxDatetime := time.Now().Local()
	inputOrderStatus := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime}

	t.Run("Success Add OrderItem", func(t *testing.T) {
		inputOrderStatusQry := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime}
		dataOrder.On("AddOrderStatuses", inputOrderStatus).Return(nil).Once()
		err := ordSrv.AddOrderStatuses(inputOrderStatusQry)

		assert.Nil(t, err)
		assert.Equal(t, inputOrderStatusQry, inputOrderStatus)
		// assert.Equal(t, inputOrderItemQry.Quantity, quantity)
	})

	t.Run("Failed Add OrderItem", func(t *testing.T) {
		expectedErr := errors.New("error add order status")

		inputOrderStatusQry := order.OrderStatusEntity{UserID: 1, Status: "", TotalPrice: 20000, PaymentUrl: paymentUrl, TrxDates: trxDatetime}
		dataOrder.On("AddOrderStatuses", inputOrderStatus).Return(expectedErr).Once()
		err := ordSrv.AddOrderStatuses(inputOrderStatusQry)

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetOrderQtyProduct(t *testing.T) {
	dataOrder := mocks.NewQueryOrderInterface(t)
	midl := mocks.NewMiddlewaresInterface(t)
	ordSrv := service.NewOrderService(dataOrder, midl)

	userID := 1
	var qty uint
	t.Run("Success GetOrderQtyProduct", func(t *testing.T) {
		dataOrder.On("GetOrderQtyProduct", uint(userID)).Return(qty, nil).Once()
		qtyprod, err := ordSrv.GetOrderQtyProduct(uint(userID))

		assert.Nil(t, err)
		assert.Equal(t, qtyprod, qty)
	})

	t.Run("Failed GetOrderQtyProduct", func(t *testing.T) {
		expectedErr := errors.New("data qty tidak ditemukan / tidak valid")

		dataOrder.On("GetOrderQtyProduct", uint(userID)).Return(qty, expectedErr).Once()
		qtyprod, err := ordSrv.GetOrderQtyProduct(uint(userID))

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, qtyprod, qty)
	})
}
