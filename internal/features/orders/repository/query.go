package repository

import (
	"fmt"
	"log"
	order "pinjamtani_project/internal/features/orders"
	"time"

	"gorm.io/gorm"
)

type orderQuery struct {
	db *gorm.DB
}

func NewOrderQuery(dbQuery *gorm.DB) order.QueryOrderInterface {
	return &orderQuery{
		db: dbQuery,
	}

}

func (q *orderQuery) AddOrderItems(orderItemData order.OrderItemEntity, orderStatusData order.OrderStatusEntity) (uint, uint, error) {

	orderItem := OrderItems{}
	query := `select c.user_id, c.product_id, c.quantity, c.id
				from "ecommerce"."carts" c
				where c.user_id = ? and c.id = ? and c."deleted_at" IS NULL;`
	cartData := Carts{}
	err := q.db.Debug().Raw(query, &orderStatusData.UserID, &orderItemData.CartID).Scan(&cartData).Error
	if err != nil {
		log.Println("Error get order item", err.Error())
		return 0, 0, err
	}
	orderItem.UserID = orderItemData.UserID
	orderItem.CartID = orderItemData.CartID
	orderItem.ProductID = cartData.ProductID
	orderItem.Quantity = cartData.Quantity

	fmt.Printf("product %d cart %d", orderItem.ProductID, orderItem.CartID)
	err = q.db.Debug().Create(&orderItem).Scan(&orderItem).Error

	if err != nil {
		log.Println("Error Insert Order", err.Error())
		return 0, 0, err
	}

	query = `UPDATE "ecommerce"."carts" SET "deleted_at" = ?, "updated_at"= ? WHERE "id"=? and "deleted_at" IS NULL;
`
	cartData = Carts{}
	// cartData.DeletedAt = time.Now()
	cartData.UpdatedAt = time.Now()
	err = q.db.Debug().Raw(query, &cartData.UpdatedAt, &cartData.UpdatedAt, &orderItemData.CartID).Scan(&cartData).Error
	if err != nil {
		log.Println("error delete cart", err.Error())

		return 0, 0, err
	}

	return orderItem.ProductID, orderItem.Quantity, nil
}

func (q *orderQuery) AddOrderStatuses(orderData order.OrderStatusEntity) error {
	// result := OrderStatus{}
	orderStatus := toOrderStatusQuery(orderData)
	err := q.db.Create(&orderStatus).Error
	if err != nil {
		log.Println("Error Insert Order", err.Error())
		return err
	}

	return nil
}
func (q *orderQuery) GetOrderQtyProduct(id uint) (uint, error) {
	var dataOrder OrderItems
	err := q.db.First(&dataOrder, id).Error
	if err != nil {
		log.Println("error get quantity product from data order item")
		return 0, err
	}
	return dataOrder.Quantity, nil

}

func (q *orderQuery) GetOrders(userid uint) (order.ListOrderItemEntity, error) {
	result := ListOrderItem{}
	query := `select
			p.id as product_id,
			p.product_name,
			p.price,
			oi.quantity, 			
			os.status,
			os.trx_dates 
			from "ecommerce".order_items oi
			join "ecommerce".products p on p.id = oi.product_id
			join "ecommerce".order_statuses os on os.user_id = oi.user_id 
			where oi.user_id = 1 and oi."deleted_at" IS NULL;`
	err := q.db.Debug().Raw(query).Scan(&result).Error
	if err != nil {
		log.Println("Error get orders", err.Error())
		return order.ListOrderItemEntity{}, err
	}
	return toListOrderItemEntity(result), nil
}
