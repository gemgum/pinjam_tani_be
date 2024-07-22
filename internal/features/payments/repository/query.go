package datapayment

import (
	"fmt"
	"log"
	payment "pinjamtani_project/internal/features/payments"
	"time"

	"gorm.io/gorm"
)

type paymentQuery struct {
	db *gorm.DB
}

func NewPaymentQuery(dbQuery *gorm.DB) payment.QueryPaymentInterface {
	return &paymentQuery{
		db: dbQuery,
	}

}

func (q *paymentQuery) UpdateStatusOrder(dataUserMidtrans payment.MidtransRequestEntity) error {
	// layout := "2006-01-02 15:04:05"

	fmt.Println("Trx Type", dataUserMidtrans.TrxType)
	fmt.Println("TrxTime ", dataUserMidtrans.TrxTime)
	fmt.Println("TrxStatus ", dataUserMidtrans.TrxStatus)
	fmt.Println("TrxID ", dataUserMidtrans.TrxID)
	fmt.Println("Message ", dataUserMidtrans.Message)
	fmt.Println("Code ", dataUserMidtrans.Code)

	datetime, err := time.Parse("2006-01-02 15:04:05", dataUserMidtrans.TrxTime)
	if err != nil {
		fmt.Println("Error conversion time:", err.Error())
		return err
	}

	query := `update "ecommerce"."order_statuses" set "status" = ? , "trx_dates" = ? where "payment_url" = ?  and "order_statuses"."deleted_at" IS NULL;`
	err = q.db.Debug().Exec(query, &dataUserMidtrans.TrxStatus, datetime, &dataUserMidtrans.TrxID).Error
	if err != nil {
		log.Println("error on update data ", err.Error())
		return err
	}
	return nil
}
