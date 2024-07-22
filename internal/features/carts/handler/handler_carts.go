package handler

import (
	"errors"
	"fmt"
	"net/http"
	"pinjamtani_project/app/middlewares"
	"pinjamtani_project/internal/features/carts"
	"pinjamtani_project/internal/features/products"
	"pinjamtani_project/internal/utils/responses"

	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartService    carts.ServiceCartInterface
	productService products.ServiceProductInterface
}

func NewCartHandler(cr carts.ServiceCartInterface, pr products.ServiceProductInterface) *CartHandler {
	return &CartHandler{
		cartService:    cr,
		productService: pr,
	}
}

func (cr *CartHandler) CreateCart(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		return errors.New("user id not found in context")
	}

	newCarts := CartRequest{}
	errBind := c.Bind(&newCarts)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "error bind carts: " + errBind.Error(),
		})
	}

	product, err := cr.productService.GetById(uint(newCarts.ProductID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "product not found: " + err.Error(),
		})
	}

	inputCart := carts.CartEntity{
		UserID:     uint(userID),
		ProductID:  uint(newCarts.ProductID),
		Quantity:   uint(newCarts.Quantity),
		TotalPrice: uint(product.Price) * uint(newCarts.Quantity),
	}

	data, errInsert := cr.cartService.Create(inputCart)
	if errInsert != nil {
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "carts failed to be created: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "carts failed to be created: "+errInsert.Error(), nil))
	}

	cartResponse := CreateCartResponse{
		Id:         data.CartID,
		UserID:     data.UserID,
		ProductID:  data.ProductID,
		Quantity:   data.Quantity,
		TotalPrice: data.TotalPrice,
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "carts was successfully created", cartResponse))
}

func (cr *CartHandler) GetAllCart(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		return errors.New("user id not found in context")
	}

	result, err := cr.cartService.GetAllCart()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "failed to retrieve all carts",
		})
	}

	var allCartResponse []CartResponse
	for _, value := range result {
		product, err := cr.productService.GetById(value.ProductID)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  "failed",
				"message": "product not found: " + err.Error(),
			})
		}
		if uint(userID) == value.UserID {
			allCartResponse = append(allCartResponse, CartResponse{
				Id:          value.CartID,
				Images:      product.Images,
				ProductName: product.ProductName,
				Price:       product.Price,
				Quantity:    float64(value.Quantity),
				TotalPrice:  float64(value.TotalPrice),
			})
		}
	}
	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "all carts fetched successfully", allCartResponse))
}

func (cr *CartHandler) UpdateCart(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		return errors.New("user id not found in context")
	}

	id := c.Param("id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "error converting id: " + errConv.Error(),
		})
	}

	var updateData CartRequest
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "error binding cart: " + err.Error(),
		})
	}

	product, err := cr.productService.GetById(uint(updateData.ProductID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "product not found: " + err.Error(),
		})
	}
	totalPrice := uint(updateData.Quantity) * uint(product.Price)

	inputArtikel := carts.CartEntity{
		ProductID:  uint(updateData.ProductID),
		Quantity:   uint(updateData.Quantity),
		TotalPrice: totalPrice,
	}

	if idConv == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "invalid cart id: idConv = 0",
		})
	}
	if inputArtikel.ProductID == 0 || inputArtikel.Quantity == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "cart product_id/quantity cannot be empty",
		})
	}

	cartData, err := cr.cartService.FindById(uint(idConv))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "cart not found: " + err.Error(),
		})
	}

	if int(cartData.UserID) != userID {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "user id not match, cannot update carts: cart user id = 0",
		})
	}

	if errInsert := cr.cartService.Update(uint(idConv), inputArtikel); errInsert != nil {
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "failed to update the cart: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to update the cart: "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "successfully updated the cart", nil))
}

func (cr *CartHandler) DeleteCart(c echo.Context) error {
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		return errors.New("user id not found in context")
	}

	id := c.Param("id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		return errors.New("error convert id")
	}

	if idConv == 0 {
		return errors.New("invalid cart id")
	}
	cartData, err := cr.cartService.FindById(uint(idConv))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "cart not found: " + err.Error(),
		})
	}

	if int(cartData.UserID) != userID {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "user id not match, cannot update carts: cart user id not founc " + fmt.Sprintf("%d", userID),
		})
	}
	if errInsert := cr.cartService.Delete(uint(idConv)); errInsert != nil {
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "Failed to delete the carts: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to delete the carts: "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "successfully deleted a carts", nil))
}
