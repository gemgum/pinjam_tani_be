package service

import (
	"log"
	payment "pinjamtani_project/internal/features/payments"
)

type paymentservice struct {
	qry payment.QueryPaymentInterface
}

func NewPaymentService(q payment.QueryPaymentInterface) payment.ServicePaymentIinterface {
	return &paymentservice{
		qry: q,
	}
}

func (p *paymentservice) UpdatePaymentStatus(reqMidtrans payment.MidtransRequestEntity) error {
	err := p.qry.UpdateStatusOrder(reqMidtrans)
	if err != nil {
		log.Println("Error on Service", err.Error())
		return err
	}
	return nil

}
