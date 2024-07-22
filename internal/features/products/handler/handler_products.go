package handler

import (
	"errors"
	"log"
	"net/http"
	"pinjamtani_project/app/middlewares"
	"pinjamtani_project/internal/features/products"
	"pinjamtani_project/internal/features/users"
	"pinjamtani_project/internal/utils/responses"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService products.ServiceProductInterface
	userService    users.ServiceUserInterface
}

func New(pr products.ServiceProductInterface, us users.ServiceUserInterface) *ProductHandler {
	return &ProductHandler{
		productService: pr,
		userService:    us,
	}
}

func (pr *ProductHandler) CreateProduct(c echo.Context) error {
	// Extract user ID from authentication context
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		err := errors.New("user id not found in context")
		log.Println(err)
		return err
	}

	// Read data from request body
	newProducts := ProductRequest{}
	errBind := c.Bind(&newProducts)
	if errBind != nil {
		err := errors.New("error bind products: " + errBind.Error())
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	// Read user image file (if any)
	file, err := c.FormFile("images")
	var imageURL string
	if err == nil {
		// Open file
		src, err := file.Open()
		if err != nil {
			errMsg := "failed to open image file: " + err.Error()
			log.Println(errMsg)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", errMsg, nil))
		}
		defer src.Close()

		// Upload file to Cloudinary
		imageURL, err = newProducts.uploadToCloudinary(src, file.Filename)
		if err != nil {
			errMsg := "failed to upload image: " + err.Error()
			log.Println(errMsg)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", errMsg, nil))
		}
	}

	// Mapping request to products
	inputProduct := products.ProductsEntity{
		UserID:      uint(userID),
		Images:      imageURL,
		ProductName: newProducts.ProductName,
		Price:       newProducts.Price,
		Stock:       newProducts.Stock,
		City:        newProducts.City,
		Description: newProducts.Description,
	}

	productID, errInsert := pr.productService.Create(inputProduct)
	if errInsert != nil {
		log.Println("Error creating product:", errInsert)
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "products failed to be created: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "products failed to be created: "+errInsert.Error(), nil))
	}

	response := ProductResponse{
		ProductID:   productID,
		Images:      imageURL,
		ProductName: newProducts.ProductName,
		Price:       newProducts.Price,
		Stock:       newProducts.Stock,
		City:        newProducts.City,
		Description: newProducts.Description,
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "products was successfully created", response))
}

func (pr *ProductHandler) GetAllProduct(c echo.Context) error {
	productName := c.QueryParam("product")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = 10
	}

	productList, err := pr.productService.GetAllProduct(productName, page, pageSize)
	if err != nil {
		log.Println("Error retrieving all products:", err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "failed to retrieve all products",
		})
	}
	var allProductResponse []ProductResponse
	for _, value := range productList {
		allProductResponse = append(allProductResponse, ProductResponse{
			ProductID:   value.ProductID,
			Images:      value.Images,
			ProductName: value.ProductName,
			Price:       value.Price,
			Stock:       value.Stock,
			City:        value.City,
		})
	}
	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "all products fetched successfully", allProductResponse))
}

func (pr *ProductHandler) DeleteProduct(c echo.Context) error {
	// Extract user ID from authentication context
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		err := errors.New("user id not found in context")
		log.Println(err)
		return err
	}

	id := c.Param("id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		err := errors.New("error convert id: " + errConv.Error())
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
	}
	if errInsert := pr.productService.Delete(uint(idConv), uint(userID)); errInsert != nil {
		log.Println("Error deleting product:", errInsert)
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "failed to delete the products: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to delete the products: "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "successfully deleted a products", nil))
}

func (pr *ProductHandler) UpdateProduct(c echo.Context) error {
	// Extract user ID from authentication context
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		err := errors.New("user id not found in context")
		log.Println(err)
		return err
	}

	id := c.Param("id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		err := errors.New("error converting id: " + errConv.Error())
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	var updateData ProductRequest
	if err := c.Bind(&updateData); err != nil {
		errMsg := "error binding product: " + err.Error()
		log.Println(errMsg)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": errMsg,
		})
	}

	// Read user image file (if any)
	file, err := c.FormFile("images")
	var imageURL string
	if err == nil {
		// Open file
		src, err := file.Open()
		if err != nil {
			errMsg := "failed to open image file: " + err.Error()
			log.Println(errMsg)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", errMsg, nil))
		}
		defer src.Close()

		// Upload file to Cloudinary
		imageURL, err = updateData.uploadToCloudinary(src, file.Filename)
		if err != nil {
			errMsg := "failed to upload image: " + err.Error()
			log.Println(errMsg)
			return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "error", errMsg, nil))
		}
	}

	inputProduct := products.ProductsEntity{
		Images:      imageURL,
		ProductName: updateData.ProductName,
		Price:       updateData.Price,
		Stock:       updateData.Stock,
		City:        updateData.City,
		Description: updateData.Description,
	}

	if errInsert := pr.productService.Update(uint(idConv), uint(userID), inputProduct); errInsert != nil {
		log.Println("Error updating product:", errInsert)
		if strings.Contains(errInsert.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, responses.JSONWebResponse(http.StatusBadRequest, "failed", "failed to update the product: "+errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.JSONWebResponse(http.StatusInternalServerError, "failed", "failed to update the product: "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusCreated, responses.JSONWebResponse(http.StatusCreated, "success", "successfully updated the product", nil))
}

func (pr *ProductHandler) GetProductId(c echo.Context) error {
	id := c.Param("id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		err := errors.New("error converting id: " + errConv.Error())
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	result, err := pr.productService.GetById(uint(idConv))
	if err != nil {
		log.Println("Error retrieving product:", err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "failed to retrieve product: " + err.Error(),
		})
	}

	userProfile, err := pr.userService.GetProfile(result.UserID)
	if err != nil {
		log.Println("Error retrieving user profile:", err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "failed to retrieve user profile: " + err.Error(),
		})
	}

	productResponse := ProductResponse{
		ProductID:    result.ProductID,
		SellerImages: userProfile.Images,
		UserName:     userProfile.UserName,
		Images:       result.Images,
		ProductName:  result.ProductName,
		Price:        result.Price,
		Stock:        result.Stock,
		City:         result.City,
		Description:  result.Description,
	}

	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "product fetched successfully", productResponse))
}

func (pr *ProductHandler) GetUserProducts(c echo.Context) error {
	// Extract user ID from authentication context
	userID := middlewares.NewMiddlewares().ExtractTokenUserId(c)
	if userID == 0 {
		err := errors.New("user id not found in context")
		log.Println(err)
		return err
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = 10
	}

	userProductList, err := pr.productService.GetUserProducts(uint(userID), page, pageSize)
	if err != nil {
		log.Println("Error retrieving user products:", err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "failed",
			"message": "failed to retrieve user products",
		})
	}
	var userProductResponse []ProductResponse
	for _, value := range userProductList {
		userProductResponse = append(userProductResponse, ProductResponse{
			ProductID:   value.ProductID,
			Images:      value.Images,
			ProductName: value.ProductName,
			Price:       value.Price,
			Stock:       value.Stock,
			Description: value.Description,
		})
	}
	return c.JSON(http.StatusOK, responses.JSONWebResponse(http.StatusOK, "success", "user products fetched successfully", userProductResponse))
}
