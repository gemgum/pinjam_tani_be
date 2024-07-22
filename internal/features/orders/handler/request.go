package handler

import (
	"fmt"
	order "projectBE23/internal/features/orders"
)

type OrderItems struct {
	CartID    uint `json:"cart_id"`
	Quantity  uint `json:"quantity"`
	ProductID uint `json:"product_id"`
	OrderID   uint `json:"order_id"`
	UserID    uint
}

type OrderStatusRequest struct {
	OrderID    uint `json:"order_id"`
	UserID     uint
	TotalPrice uint   `json:"total_price"`
	Email      string `json:"email"`
	Phone      uint   `json:"phone"`
	FullName   string `json:"fullname"`
	Cart       []uint `json:"cart_id"`
}

type ResponseMidtransData struct {
	Token         string
	RedirectUrl   string
	StatusCode    string
	ErrorMessages []string
}

type RequestMidtransData struct {
	Fullname string
	Email    string
	Phone    string
	OrderID  string
	GrossAmt int64
}

func toMidtransData(dataOrderMidtrans OrderStatusRequest, orderID string) order.MidatransDataRequestEntity {
	return order.MidatransDataRequestEntity{
		Fullname: dataOrderMidtrans.FullName,
		Email:    dataOrderMidtrans.Email,
		Phone:    fmt.Sprintf("%d", dataOrderMidtrans.Phone),
		GrossAmt: int64(dataOrderMidtrans.TotalPrice),
		OrderID:  orderID,
	}
}

func toOrderItemRequest(orderItemData OrderItems) order.OrderItemEntity {
	return order.OrderItemEntity{
		UserID:    orderItemData.UserID,
		ProductID: orderItemData.ProductID,
		Quantity:  orderItemData.Quantity,
		CartID:    orderItemData.CartID,
	}

}

func toOrderStatusRequest(orderStatusData OrderStatusRequest, orderID string) order.OrderStatusEntity {
	return order.OrderStatusEntity{
		TotalPrice: orderStatusData.TotalPrice,
		UserID:     orderStatusData.UserID,
		PaymentUrl: orderID,
	}

}
