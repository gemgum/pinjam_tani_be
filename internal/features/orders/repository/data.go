package repository

import (
	"fmt"
	order "projectBE23/internal/features/orders"
	"time"

	"gorm.io/gorm"
)

type OrderItems struct {
	gorm.Model
	UserID    uint
	CartID    uint
	ProductID uint
	Quantity  uint
}

type Carts struct {
	gorm.Model
	CartID    uint `json:"id"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type OrderStatus struct {
	gorm.Model
	UserID     uint
	Status     string
	TotalPrice uint
	PaymentUrl string
	TrxDates   time.Time
}

type ListOrderItem []struct {
	ProductID   uint64    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       uint      `json:"price"`
	Quantity    uint      `json:"quantity"`
	TrxDate     time.Time `json:"trx_dates"`
	Status      string    `json:"status"`
}

func toOrderStatusQuery(orderStatusData order.OrderStatusEntity) OrderStatus {
	return OrderStatus{
		// PaymentUrl: orderStatusData.PaymentUrl,
		// UserID:     orderStatusData.UserID,
		TotalPrice: orderStatusData.TotalPrice,
		UserID:     orderStatusData.UserID,
		PaymentUrl: orderStatusData.PaymentUrl,
	}

}

func toListOrderItemEntity(data ListOrderItem) order.ListOrderItemEntity {
	var result order.ListOrderItemEntity
	for _, v := range data {
		fmt.Printf("id Q %d\n", v.ProductID)
		dataList := struct {
			ProductID   uint64
			ProductName string
			Price       uint
			Quantity    uint
			TrxDate     time.Time
			Status      string
		}{
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			Price:       v.Price,
			Quantity:    v.Quantity,
			TrxDate:     v.TrxDate,
			Status:      v.Status,
		}
		result = append(result, dataList)
	}

	return result
}
