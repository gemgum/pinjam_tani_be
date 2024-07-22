package handler

import (
	"fmt"
	"log"
	payment "projectBE23/internal/features/payments"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"strings"
)

type PaymentRequestMidtrans struct {
	TrxType   string `json:"payment_type"`
	TrxTime   string `json:"transaction_time"`
	TrxStatus string `json:"transaction_status"`
	TrxID     string `json:"order_id"`
	Message   string `json:"status_message"`
	Code      string `json:"status_code"`
}

type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Field   string      `json:"field"`
	Message string      `json:"message"`
}

func getFName(name payment.PaymentEntity) string {
	dataN := strings.Index(name.Fullname, " ")
	splitResult := strings.SplitN(name.Fullname, " ", dataN)
	rv := fmt.Sprintf("%q", splitResult)
	return rv
}

func getLName(name payment.PaymentEntity) string {
	// dataN := strings.Index(name.Fullname, " ")
	splitResult := strings.SplitAfter(name.Fullname, " ")
	rv := fmt.Sprintf("%q", splitResult)
	return rv
}

func ToMidatransRequest(dataPayment payment.PaymentEntity) *snap.Request {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{

			OrderID:  dataPayment.OrderID,
			GrossAmt: dataPayment.GrossAmt,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: getFName(dataPayment),
			LName: getLName(dataPayment),
			Email: dataPayment.Email,
			Phone: dataPayment.Phone,
		},
	}
	log.Println("Order ", req.TransactionDetails.OrderID)
	log.Println("Gross amnt", req.TransactionDetails.GrossAmt)

	log.Println("Fname ", req.CustomerDetail.FName)

	log.Println("Lname", req.CustomerDetail.LName)

	log.Println("Email ", req.CustomerDetail.Email)

	log.Println("Phone ", req.CustomerDetail.Phone)

	return req
}

func toMidtransEntity(data snap.Response) payment.MidtransData {
	rv := payment.MidtransData{
		RedirectUrl: data.RedirectURL,
		Token:       data.Token,
	}

	return rv
}
