package order

import (
	"time"

	"github.com/labstack/echo/v4"
)

type OrderItemEntity struct {
	ProductID uint
	CartID    uint
	Quantity  uint
	UserID    uint
}

type OrderStatusEntity struct {
	UserID     uint
	Status     string
	TotalPrice uint
	PaymentUrl string
	TrxDates   time.Time
}

type MidatransDataRequestEntity struct {
	Fullname string
	Email    string
	Phone    string
	OrderID  string
	GrossAmt int64
}

type MidtransDataEntity struct {
	Token         string
	RedirectUrl   string
	StatusCode    string
	ErrorMessages []string
}

type ListOrderItemEntity []struct {
	ProductID   uint64
	ProductName string
	Price       uint
	Quantity    uint
	TrxDate     time.Time
	Status      string
}

type MidtransRequestEntity struct {
	TrxType   string
	TrxTime   string
	TrxStatus string
	TrxID     string
	Message   string
	Code      string
}

type QueryOrderInterface interface {
	AddOrderItems(orderItemData OrderItemEntity, orderStatusData OrderStatusEntity) (uint, uint, error)
	AddOrderStatuses(orderData OrderStatusEntity) error
	GetOrders(userid uint) (ListOrderItemEntity, error)
	GetOrderQtyProduct(id uint) (uint, error)
}

type HandlerOrderInterface interface {
	AddOrderItems() echo.HandlerFunc
	GetOrders() echo.HandlerFunc
}

type ServiceOrderInterface interface {
	AddOrderItems(orderItemData OrderItemEntity, orderStatusData OrderStatusEntity) (uint, uint, error)
	AddOrderStatuses(orderData OrderStatusEntity) error
	ProccesOrderPayment(userMidtransData MidatransDataRequestEntity) (MidtransDataEntity, error)
	GetOrders(userID uint) (ListOrderItemEntity, error)
	GetOrderQtyProduct(id uint) (uint, error)
	// DoPrccesPayment(dataOrder OrderStatusEntity) (MidtransDataEntity, error)
}
