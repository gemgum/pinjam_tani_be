package payment

import (
	"github.com/labstack/echo/v4"
)

type PaymentEntity struct {
	Fullname string
	Email    string
	Phone    string
	OrderID  string
	GrossAmt int64
}

type MidtransData struct {
	Token         string
	RedirectUrl   string
	StatusCode    string
	ErrorMessages []string
}

type MidtransRequestEntity struct {
	TrxType   string
	TrxTime   string
	TrxStatus string
	TrxID     string
	Message   string
	Code      string
}

type HandlerPaymentInterface interface {
	ProcessPayment(dataPayment PaymentEntity) (MidtransData, error)
	UpdateStatusOrderMidtrans() echo.HandlerFunc
}

type ServicePaymentIinterface interface {
	// ProcessPayment(dataPayment MidtransData) error
	UpdatePaymentStatus(reqMidtrans MidtransRequestEntity) error
}

type QueryPaymentInterface interface {
	UpdateStatusOrder(dataUserMidtrans MidtransRequestEntity) error
}
