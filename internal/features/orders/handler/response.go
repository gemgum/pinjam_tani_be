package handler

import (
	"fmt"
	order "pinjamtani_project/internal/features/orders"
	"time"
)

type OrderResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	StatusCode    string   `json:"status_code,omitempty"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

type ListOrderItemResponse []struct {
	ProducID    uint64    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       uint      `json:"price"`
	Quantity    uint      `json:"quantity"`
	TrxDate     time.Time `json:"trx_dates"`
	Status      string    `json:"status"`
}

func toOrderResponse(dataOrderMidtrans order.MidtransDataEntity) OrderResponse {
	return OrderResponse{
		Token:         dataOrderMidtrans.Token,
		RedirectURL:   dataOrderMidtrans.RedirectUrl,
		StatusCode:    dataOrderMidtrans.StatusCode,
		ErrorMessages: dataOrderMidtrans.ErrorMessages,
	}
}

func toResponseFormat(listData order.ListOrderItemEntity) ListOrderItemResponse {
	var result ListOrderItemResponse
	for _, v := range listData {
		fmt.Printf("id S  %d\n", v.ProductID)
		dataList := struct {
			ProducID    uint64    `json:"product_id"`
			ProductName string    `json:"product_name"`
			Price       uint      `json:"price"`
			Quantity    uint      `json:"quantity"`
			TrxDate     time.Time `json:"trx_dates"`
			Status      string    `json:"status"`
		}{
			ProducID:    v.ProductID,
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
