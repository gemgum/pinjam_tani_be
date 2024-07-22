package handler

import (
	"log"
	"net/http"
	"projectBE23/app/config"
	payment "projectBE23/internal/features/payments"
	"projectBE23/internal/helper"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentHandler struct {
	snap *snap.Client
	srv  payment.ServicePaymentIinterface
}

func NewPaymentHandler(s *snap.Client, sr payment.ServicePaymentIinterface) payment.HandlerPaymentInterface {
	return &paymentHandler{
		snap: s,
		srv:  sr,
	}
}
func (ph *paymentHandler) ProcessPayment(dataPayment payment.PaymentEntity) (payment.MidtransData, error) {
	ph.snap.New(config.SERVER_KEY, midtrans.Sandbox)
	snapResp, err := ph.snap.CreateTransaction(ToMidatransRequest(dataPayment))

	if err != nil {
		log.Panicln("Error From Midtrans", err.RawError)
		return payment.MidtransData{}, err.RawError
	}
	return toMidtransEntity(*snapResp), nil
}

func (ph *paymentHandler) UpdateStatusOrderMidtrans() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input PaymentRequestMidtrans
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bad request", err.Error())
			return c.JSON(http.StatusBadRequest, http.ErrBodyNotAllowed)
		}
		err = ph.srv.UpdatePaymentStatus(payment.MidtransRequestEntity(input))
		if err != nil {
			log.Println("error on handler")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on server", nil))

		}
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "success", nil))

	}
}
