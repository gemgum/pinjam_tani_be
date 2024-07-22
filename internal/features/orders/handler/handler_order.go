package handler

import (
	"fmt"
	"log"
	"net/http"
	"projectBE23/app/middlewares"
	order "projectBE23/internal/features/orders"
	products "projectBE23/internal/features/products"
	"projectBE23/internal/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	srv  order.ServiceOrderInterface
	pQry products.DataProductInterface
}

func NewOrderHandler(s order.ServiceOrderInterface, p products.DataProductInterface) order.HandlerOrderInterface {
	return &orderHandler{
		srv:  s,
		pQry: p,
	}
}
func (h *orderHandler) GetOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
		if userID < 1 {
			log.Println("Error Unautorized ", userID)
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Error Unatorized", nil))
		}
		result, err := h.srv.GetOrders(uint(userID))
		if err != nil {
			log.Println("Error on server ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error in server", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseGetOrderFormat(http.StatusOK, "success", userID, toResponseFormat(result)))
	}
}

func (h *orderHandler) AddOrderItems() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input OrderStatusRequest
		userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
		if userID < 1 {
			log.Println("Error Unautorized ", userID)
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Error Unatorized", nil))
		}
		input.UserID = uint(userID)
		err := c.Bind(&input)
		if err != nil {
			log.Println("Error on data ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "error input data", nil))
		}

		var ordetItem OrderItems
		for _, val := range input.Cart {
			ordetItem.UserID = uint(userID)
			ordetItem.CartID = val
			fmt.Printf("Cart :  %d \n", val)
			prodID, qty, err := h.srv.AddOrderItems(toOrderItemRequest(ordetItem), toOrderStatusRequest(input, ""))
			if err != nil {
				log.Println("Error on service ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
			}
			err = h.pQry.DecreaseProduct(prodID, qty)
			if err != nil {
				log.Println("Error update stock product ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
			}
		}
		orderID := uuid.NewString()
		err = h.srv.AddOrderStatuses(toOrderStatusRequest(input, orderID))
		if err != nil {
			log.Println("Error on service ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error on database", nil))
		}

		dataMid, err := h.srv.ProccesOrderPayment(toMidtransData(input, orderID))
		if err != nil {
			log.Println("Erro when hit midtransdata ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "error in server", nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "success", toOrderResponse(dataMid)))
	}
}
func (h *orderHandler) UpdateOrderStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
