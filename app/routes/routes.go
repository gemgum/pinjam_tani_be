package routes

import (
	"projectBE23/app/middlewares"
	"projectBE23/internal/utils"
	"projectBE23/internal/utils/encrypts"

	_productHandler "projectBE23/internal/features/products/handler"
	_productData "projectBE23/internal/features/products/repository"
	_productService "projectBE23/internal/features/products/service"
	_userHandler "projectBE23/internal/features/users/handler"
	_userData "projectBE23/internal/features/users/repository"
	_userService "projectBE23/internal/features/users/service"

	orderHandler "projectBE23/internal/features/orders/handler"
	orderData "projectBE23/internal/features/orders/repository"
	orderService "projectBE23/internal/features/orders/service"

	paymentHandler "projectBE23/internal/features/payments/handler"
	paymentData "projectBE23/internal/features/payments/repository"
	paymentService "projectBE23/internal/features/payments/service"

	cartHandler "projectBE23/internal/features/carts/handler"
	cartData "projectBE23/internal/features/carts/repository"
	cartService "projectBE23/internal/features/carts/service"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB) {
	middlewares := middlewares.NewMiddlewares()
	hashService := encrypts.NewHashService()
	accountUtility := utils.NewAccountUtility()

	userData := _userData.New(db)
	userService := _userService.New(userData, hashService, middlewares, accountUtility)
	userHandlerAPI := _userHandler.New(userService)

	productData := _productData.NewPaymentQuery(db)
	productService := _productService.New(productData)
	productHandlerAPI := _productHandler.New(productService, userService)

	cartData := cartData.New(db)
	cartSrv := cartService.New(cartData)
	cartHandler := cartHandler.New(cartSrv, productService)

	ordData := orderData.NewOrderQuery(db)
	ordSrv := orderService.NewOrderService(ordData, middlewares)
	ordHandler := orderHandler.NewOrderHandler(ordSrv, productData)

	//userHandler
	e.POST("/users", userHandlerAPI.Register)
	e.POST("/login", userHandlerAPI.Login)
	e.POST("/logout", userHandlerAPI.Logout, middlewares.JWTMiddleware())
	e.GET("/users", userHandlerAPI.GetProfile, middlewares.JWTMiddleware())
	e.PUT("/users", userHandlerAPI.Update, middlewares.JWTMiddleware())
	e.DELETE("/users", userHandlerAPI.Delete, middlewares.JWTMiddleware())

	s := snap.Client{}
	q := paymentData.NewPaymentQuery(db)
	sPay := paymentService.NewPaymentService(q)
	pay := paymentHandler.NewPaymentHandler(&s, sPay)

	e.POST("/update_payment", pay.UpdateStatusOrderMidtrans())
	a := e.Group("/orders")

	a.POST("", ordHandler.AddOrderItems())
	a.GET("", ordHandler.GetOrders())

	//productHandler
	e.POST("/products", productHandlerAPI.CreateProduct, middlewares.JWTMiddleware())
	e.GET("/products", productHandlerAPI.GetAllProduct)
	e.GET("/users/products", productHandlerAPI.GetUserProducts, middlewares.JWTMiddleware())
	e.GET("/products/:id", productHandlerAPI.GetProductId, middlewares.JWTMiddleware())
	e.PUT("/products/:id", productHandlerAPI.UpdateProduct, middlewares.JWTMiddleware())
	e.DELETE("/products/:id", productHandlerAPI.DeleteProduct, middlewares.JWTMiddleware())

	b := e.Group("/carts")
	b.POST("", cartHandler.CreateCart, middlewares.JWTMiddleware())
	b.GET("", cartHandler.GetAllCart, middlewares.JWTMiddleware())
	b.PUT("/:id", cartHandler.UpdateCart, middlewares.JWTMiddleware())
	b.DELETE("/:id", cartHandler.DeleteCart, middlewares.JWTMiddleware())
}
