package service

import (
	"errors"
	"fmt"
	"log"
	"pinjamtani_project/app/middlewares"
	order "pinjamtani_project/internal/features/orders"
	payment "pinjamtani_project/internal/features/payments"
	"pinjamtani_project/internal/features/payments/handler"
	datapayment "pinjamtani_project/internal/features/payments/repository"
	"pinjamtani_project/internal/features/payments/service"

	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type orderServices struct {
	midl middlewares.MiddlewaresInterface
	qry  order.QueryOrderInterface
}

func NewOrderService(q order.QueryOrderInterface, m middlewares.MiddlewaresInterface) order.ServiceOrderInterface {
	return &orderServices{
		midl: m,
		qry:  q,
	}
}

func (o *orderServices) AddOrderItems(orderItemData order.OrderItemEntity, orderStatus order.OrderStatusEntity) (uint, uint, error) {
	fmt.Println("chart item", orderItemData.CartID)
	productID, qty, err := o.qry.AddOrderItems(orderItemData, orderStatus)
	if err != nil {
		log.Println("Error Service", err.Error())
		return 0, 0, err
	}
	return productID, qty, nil
}

func (o *orderServices) AddOrderStatuses(orderData order.OrderStatusEntity) error {
	err := o.qry.AddOrderStatuses(orderData)
	if err != nil {
		log.Println("Error Service", err.Error())
		return err
	}
	return nil
}

func (o *orderServices) GetOrders(userID uint) (order.ListOrderItemEntity, error) {
	result, err := o.qry.GetOrders(userID)
	if err != nil {
		log.Println("Error Service", err.Error())
		return order.ListOrderItemEntity{}, errors.New("data id tidak ditemukan / tidak valid")
	}
	return result, nil
}

func (o *orderServices) GetOrderQtyProduct(id uint) (uint, error) {
	result, err := o.qry.GetOrderQtyProduct(id)
	if err != nil {
		log.Println("Error Service", err.Error())
		return 0, err
	}
	return result, nil
}

func (o *orderServices) ProccesOrderPayment(dataMidtrans order.MidatransDataRequestEntity) (order.MidtransDataEntity, error) {
	s := snap.Client{}
	q := datapayment.NewPaymentQuery(&gorm.DB{})
	sPay := service.NewPaymentService(q)
	pay := handler.NewPaymentHandler(&s, sPay)
	dataMid, err := pay.ProcessPayment(ToPaymentOrderEntity(dataMidtrans))
	if err != nil {
		log.Println("Error Service", err.Error())
		return toMidtransDataEntity(dataMid), err
	}

	return toMidtransDataEntity(dataMid), nil
}

// func (o *orderServices) DoPrccesPayment(dataOrderStatus order.OrderStatusEntity, dataOrderItems order.OrderItemEntity) (order.MidtransDataEntity, error) {
// 	mitransData := order.MidtransDataEntity{}

// 	var ordetItem order.OrderItemEntity
// 	// for _, val := range input.Cart {
// 	// 	ordetItem.UserID = uint(userID)
// 	// 	ordetItem.CartID = val
// 	// 	fmt.Printf("Cart :  %d \n", val)
// 	// 	prodID, qty, err := h.srv.AddOrderItems(toOrderItemRequest(ordetItem), toOrderStatusRequest(input, ""))
// 	// 	if err != nil {
// 	// 		log.Println("Error on service ", err.Error())
// 	// 		return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
// 	// 	}
// 	// 	err = h.pQry.DecreaseProduct(prodID, qty)
// 	// 	if err != nil {
// 	// 		log.Println("Error update stock product ", err.Error())
// 	// 		return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
// 	// 	}
// 	// }
// 	orderID := uuid.NewString()
// 	err := o.AddOrderStatuses(dataOrderStatus)
// 	if err != nil {
// 		log.Println("Error on service ", err.Error())
// 		return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
// 	}

// 	dataMid, err := o.ProccesOrderPayment(toMidtransData(input, orderID))
// 	if err != nil {
// 		log.Println("Erro when hit midtransdata ", err.Error())
// 		return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error in server", nil))
// 	}
// 	return mitransData, nil
// }

func ToPaymentOrderEntity(dataMidtrans order.MidatransDataRequestEntity) payment.PaymentEntity {
	return payment.PaymentEntity{
		Fullname: dataMidtrans.Fullname,
		Email:    dataMidtrans.Email,
		Phone:    dataMidtrans.Phone,
		OrderID:  dataMidtrans.OrderID,
		GrossAmt: dataMidtrans.GrossAmt,
	}

}

func toMidtransDataEntity(dataMidtrans payment.MidtransData) order.MidtransDataEntity {
	return order.MidtransDataEntity{
		Token:         dataMidtrans.Token,
		RedirectUrl:   dataMidtrans.RedirectUrl,
		ErrorMessages: dataMidtrans.ErrorMessages,
		StatusCode:    dataMidtrans.StatusCode,
	}

}
