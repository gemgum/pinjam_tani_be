package routes

import (
	"pinjamtani_project/app/middlewares"
	"pinjamtani_project/internal/utils"
	"pinjamtani_project/internal/utils/encrypts"

	_productHandler "pinjamtani_project/internal/features/products/handler"
	_pD "pinjamtani_project/internal/features/products/repository"
	_pS "pinjamtani_project/internal/features/products/service"

	_userHandler "pinjamtani_project/internal/features/users/handler"
	_uD "pinjamtani_project/internal/features/users/repository"
	_uS "pinjamtani_project/internal/features/users/service"

	orderHandler "pinjamtani_project/internal/features/orders/handler"
	orderData "pinjamtani_project/internal/features/orders/repository"
	orderService "pinjamtani_project/internal/features/orders/service"

	paymentHandler "pinjamtani_project/internal/features/payments/handler"
	paymentData "pinjamtani_project/internal/features/payments/repository"
	paymentService "pinjamtani_project/internal/features/payments/service"

	cartHandler "pinjamtani_project/internal/features/carts/handler"
	cartData "pinjamtani_project/internal/features/carts/repository"
	cartService "pinjamtani_project/internal/features/carts/service"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB) {
	middlewares := middlewares.NewMiddlewares()
	hashService := encrypts.NewHashService()
	accountUtility := utils.NewAccountUtility()

	uD := _uD.NewUserService(db)
	uS := _uS.NewUserService(uD, hashService, middlewares, accountUtility)
	uH := _userHandler.NewUserHandler(uS)

	pD := _pD.NewPaymentQuery(db)
	pS := _pS.New(pD)
	pH := _productHandler.New(pS, uS)

	cartData := cartData.NewCartData(db)
	cartSrv := cartService.NewCartService(cartData)
	cartHandler := cartHandler.NewCartHandler(cartSrv, pS)

	ordData := orderData.NewOrderQuery(db)
	ordSrv := orderService.NewOrderService(ordData, middlewares)
	ordHandler := orderHandler.NewOrderHandler(ordSrv, pD)

	e.POST("/users", uH.Register)
	e.POST("/login", uH.Login)
	e.POST("/logout", uH.Logout, middlewares.JWTMiddleware())
	e.GET("/users", uH.GetProfile, middlewares.JWTMiddleware())
	e.PUT("/users", uH.Update, middlewares.JWTMiddleware())
	e.DELETE("/users", uH.Delete, middlewares.JWTMiddleware())

	s := snap.Client{}
	q := paymentData.NewPaymentQuery(db)
	sPay := paymentService.NewPaymentService(q)
	pay := paymentHandler.NewPaymentHandler(&s, sPay)

	e.POST("/update_payment", pay.UpdateStatusOrderMidtrans())
	a := e.Group("/orders")

	a.POST("", ordHandler.AddOrderItems())
	a.GET("", ordHandler.GetOrders())

	e.POST("/products", pH.CreateProduct, middlewares.JWTMiddleware())
	e.GET("/products", pH.GetAllProduct)
	e.GET("/users/products", pH.GetUserProducts, middlewares.JWTMiddleware())
	e.GET("/products/:id", pH.GetProductId, middlewares.JWTMiddleware())
	e.PUT("/products/:id", pH.UpdateProduct, middlewares.JWTMiddleware())
	e.DELETE("/products/:id", pH.DeleteProduct, middlewares.JWTMiddleware())

	b := e.Group("/carts")
	b.POST("", cartHandler.CreateCart, middlewares.JWTMiddleware())
	b.GET("", cartHandler.GetAllCart, middlewares.JWTMiddleware())
	b.PUT("/:id", cartHandler.UpdateCart, middlewares.JWTMiddleware())
	b.DELETE("/:id", cartHandler.DeleteCart, middlewares.JWTMiddleware())
}
