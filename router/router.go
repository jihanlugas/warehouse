package router

import (
	"encoding/json"
	"fmt"
	"github.com/jihanlugas/warehouse/app/auth"
	"github.com/jihanlugas/warehouse/app/customer"
	"github.com/jihanlugas/warehouse/app/inbound"
	"github.com/jihanlugas/warehouse/app/outbound"
	"github.com/jihanlugas/warehouse/app/product"
	"github.com/jihanlugas/warehouse/app/purchaseorder"
	"github.com/jihanlugas/warehouse/app/purchaseorderproduct"
	"github.com/jihanlugas/warehouse/app/retail"
	"github.com/jihanlugas/warehouse/app/retailproduct"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stockin"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovement"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicle"
	"github.com/jihanlugas/warehouse/app/transaction"
	"github.com/jihanlugas/warehouse/app/user"
	"github.com/jihanlugas/warehouse/app/userprivilege"
	"github.com/jihanlugas/warehouse/app/vehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/constant"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/response"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/jihanlugas/warehouse/docs"
	"net/http"
)

func Init() *echo.Echo {
	//// repositories
	userRepository := user.NewRepository()
	userprivilegeRepository := userprivilege.NewRepository()
	warehouseRepository := warehouse.NewRepository()
	transactionRepository := transaction.NewRepository()
	vehicleRepository := vehicle.NewRepository()
	customerRepository := customer.NewRepository()
	productRepository := product.NewRepository()
	retailRepository := retail.NewRepository()
	retailproductRepository := retailproduct.NewRepository()
	purchaseorderRepository := purchaseorder.NewRepository()
	purchaseorderproductRepository := purchaseorderproduct.NewRepository()
	outboundRepository := outbound.NewRepository()
	inboundRepository := inbound.NewRepository()
	stockinRepository := stockin.NewRepository()
	stockRepository := stock.NewRepository()
	stocklogRepository := stocklog.NewRepository()
	stockmovementRepository := stockmovement.NewRepository()
	stockmovementvehicleRepository := stockmovementvehicle.NewRepository()

	// usecases
	authUsecase := auth.NewUsecase(userRepository, warehouseRepository)
	userUsecase := user.NewUsecase(userRepository, userprivilegeRepository)
	warehouseUsecase := warehouse.NewUsecase(warehouseRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository)
	stocklogUsecase := stocklog.NewUsecase(stocklogRepository)
	stockmovementvehicleUsecase := stockmovementvehicle.NewUsecase(stockmovementvehicleRepository, warehouseRepository, vehicleRepository, stockmovementRepository, stockRepository, stocklogRepository, purchaseorderRepository, retailRepository)
	vehicleUsecase := vehicle.NewUsecase(vehicleRepository)
	customerUsecase := customer.NewUsecase(customerRepository)
	productUsecase := product.NewUsecase(productRepository)
	retailUsecase := retail.NewUsecase(retailRepository, retailproductRepository, stockRepository, stocklogRepository, customerRepository, warehouseRepository, vehicleRepository)
	purchaseorderUsecase := purchaseorder.NewUsecase(purchaseorderRepository, purchaseorderproductRepository, stockRepository, stocklogRepository, customerRepository, warehouseRepository, vehicleRepository)
	outboundUsecase := outbound.NewUsecase(outboundRepository, warehouseRepository, vehicleRepository, stockRepository, stocklogRepository, stockmovementRepository, stockmovementvehicleRepository)
	inboundUsecase := inbound.NewUsecase(inboundRepository, warehouseRepository, vehicleRepository, stockRepository, stocklogRepository, stockmovementRepository, stockmovementvehicleRepository)
	stockinUsecase := stockin.NewUsecase(stockinRepository, warehouseRepository, stockRepository, stocklogRepository, stockmovementRepository)

	// handlers
	authHandler := auth.NewHandler(authUsecase)
	userHandler := user.NewHandler(userUsecase)
	warehouseHandler := warehouse.NewHandler(warehouseUsecase)
	transactionHandler := transaction.NewHandler(transactionUsecase)
	stocklogHandler := stocklog.NewHandler(stocklogUsecase)
	stockmovementvehicleHandler := stockmovementvehicle.NewHandler(stockmovementvehicleUsecase)
	vehicleHandler := vehicle.NewHandler(vehicleUsecase)
	customerHandler := customer.NewHandler(customerUsecase)
	productHandler := product.NewHandler(productUsecase)
	retailHandler := retail.NewHandler(retailUsecase)
	purchaseorderHandler := purchaseorder.NewHandler(purchaseorderUsecase)
	outboundHandler := outbound.NewHandler(outboundUsecase)
	inboundHandler := inbound.NewHandler(inboundUsecase)
	stockinHandler := stockin.NewHandler(stockinUsecase)

	router := websiteRouter()

	if config.Debug {
		router.GET("/", func(c echo.Context) error {
			return response.Success(http.StatusOK, "Welcome", nil).SendJSON(c)
		})
		router.GET("/swg/*", echoSwagger.WrapHandler)
	}

	routerAuth := router.Group("/auth")
	routerAuth.POST("/sign-in", authHandler.SignIn)
	routerAuth.GET("/sign-out", authHandler.SignOut)
	routerAuth.GET("/init", authHandler.Init, checkTokenMiddleware)
	routerAuth.GET("/refresh-token", authHandler.RefreshToken, checkTokenMiddleware)

	routerCustomer := router.Group("/customer", checkTokenMiddlewareAdmin)
	routerCustomer.GET("", customerHandler.Page)
	routerCustomer.POST("", customerHandler.Create)
	routerCustomer.GET("/:id", customerHandler.GetById)
	routerCustomer.PUT("/:id", customerHandler.Update)
	routerCustomer.DELETE("/:id", customerHandler.Delete)

	routerWarehouse := router.Group("/warehouse")
	routerWarehouse.GET("", warehouseHandler.Page, checkTokenMiddleware)
	routerWarehouse.POST("", warehouseHandler.Create, checkTokenMiddlewareAdmin)
	routerWarehouse.GET("/:id", warehouseHandler.GetById, checkTokenMiddleware)
	routerWarehouse.PUT("/:id", warehouseHandler.Update, checkTokenMiddlewareAdmin)
	routerWarehouse.DELETE("/:id", warehouseHandler.Delete, checkTokenMiddlewareAdmin)

	routerTransaction := router.Group("/transaction", checkTokenMiddlewareAdmin)
	routerTransaction.GET("", transactionHandler.Page)
	routerTransaction.POST("", transactionHandler.Create)
	routerTransaction.GET("/:id", transactionHandler.GetById)
	routerTransaction.PUT("/:id", transactionHandler.Update)
	routerTransaction.DELETE("/:id", transactionHandler.Delete)

	routerStocklog := router.Group("/stocklog")
	routerStocklog.GET("", stocklogHandler.Page, checkTokenMiddleware)
	routerStocklog.GET("/:id", stocklogHandler.GetById, checkTokenMiddleware)

	routerStockmovementvehicle := router.Group("/stockmovementvehicle")
	routerStockmovementvehicle.GET("", stockmovementvehicleHandler.Page, checkTokenMiddleware)
	routerStockmovementvehicle.GET("/:id", stockmovementvehicleHandler.GetById, checkTokenMiddleware)
	routerStockmovementvehicle.DELETE("/:id", stockmovementvehicleHandler.Delete, checkTokenMiddleware)
	routerStockmovementvehicle.GET("/:id/set-sent", stockmovementvehicleHandler.SetSent, checkTokenMiddleware)
	routerStockmovementvehicle.GET("/:id/generate-delivery-order", stockmovementvehicleHandler.GenerateDeliveryOrder, checkTokenMiddleware)

	routerStockmovementvehiclePurchaseorder := routerStockmovementvehicle.Group("/purchaseorder")
	routerStockmovementvehiclePurchaseorder.POST("", stockmovementvehicleHandler.CreatePurchaseorder, checkTokenMiddleware)
	routerStockmovementvehiclePurchaseorder.PUT("/:id", stockmovementvehicleHandler.UpdatePurchaseorder, checkTokenMiddleware)

	routerStockmovementvehicleRetail := routerStockmovementvehicle.Group("/retail")
	routerStockmovementvehicleRetail.POST("", stockmovementvehicleHandler.CreateRetail, checkTokenMiddleware)
	routerStockmovementvehicleRetail.PUT("/:id", stockmovementvehicleHandler.UpdateRetail, checkTokenMiddleware)

	routerProduct := router.Group("/product")
	routerProduct.GET("", productHandler.Page, checkTokenMiddleware)
	routerProduct.POST("", productHandler.Create, checkTokenMiddlewareAdmin)
	routerProduct.GET("/:id", productHandler.GetById, checkTokenMiddleware)
	routerProduct.PUT("/:id", productHandler.Update, checkTokenMiddlewareAdmin)
	routerProduct.DELETE("/:id", productHandler.Delete, checkTokenMiddlewareAdmin)

	routerRetail := router.Group("/retail", checkTokenMiddleware)
	routerRetail.GET("", retailHandler.Page)
	routerRetail.POST("", retailHandler.Create)
	routerRetail.GET("/:id", retailHandler.GetById)
	routerRetail.PUT("/:id", retailHandler.Update)
	//routerRetail.DELETE("/:id", retailHandler.Delete)
	routerRetail.GET("/:id/set-status-open", retailHandler.SetStatusOpen)
	routerRetail.GET("/:id/set-status-close", retailHandler.SetStatusClose)
	routerRetail.GET("/:id/generate-invoice", retailHandler.GenerateInvoice)

	routerPurchaseorder := router.Group("/purchaseorder", checkTokenMiddleware)
	routerPurchaseorder.GET("", purchaseorderHandler.Page)
	routerPurchaseorder.POST("", purchaseorderHandler.Create)
	routerPurchaseorder.GET("/:id", purchaseorderHandler.GetById)
	routerPurchaseorder.PUT("/:id", purchaseorderHandler.Update)
	//routerPurchaseorder.DELETE("/:id", purchaseorderHandler.Delete)
	routerPurchaseorder.GET("/:id/set-status-open", purchaseorderHandler.SetStatusOpen)
	routerPurchaseorder.GET("/:id/set-status-close", purchaseorderHandler.SetStatusClose)
	routerPurchaseorder.GET("/:id/generate-invoice", purchaseorderHandler.GenerateInvoice)

	routerUser := router.Group("/user")
	routerUser.GET("", userHandler.Page, checkTokenMiddlewareAdmin)
	routerUser.POST("", userHandler.Create, checkTokenMiddlewareAdmin)
	routerUser.GET("/:id", userHandler.GetById, checkTokenMiddlewareAdmin)
	routerUser.PUT("/:id", userHandler.Update, checkTokenMiddlewareAdmin)
	routerUser.POST("/:id/privilege", userHandler.UpdateUserPrivilege, checkTokenMiddlewareAdmin)
	routerUser.DELETE("/:id", userHandler.Delete, checkTokenMiddlewareAdmin)
	routerUser.POST("/change-password", userHandler.Delete, checkTokenMiddleware)

	routerVehicle := router.Group("/vehicle", checkTokenMiddleware)
	routerVehicle.GET("", vehicleHandler.Page)
	routerVehicle.POST("", vehicleHandler.Create)
	routerVehicle.GET("/:id", vehicleHandler.GetById)
	routerVehicle.PUT("/:id", vehicleHandler.Update)
	routerVehicle.DELETE("/:id", vehicleHandler.Delete)

	routerInbound := router.Group("/inbound", checkTokenMiddleware)
	routerInbound.GET("", inboundHandler.Page)
	routerInbound.GET("/:id", inboundHandler.GetById)
	routerInbound.PUT("/:id", inboundHandler.Update)
	routerInbound.GET("/:id/set-recived", inboundHandler.SetRecived)
	routerInbound.GET("/:id/generate-delivery-recipt", inboundHandler.GenerateDeliveryRecipt)

	routerOutbound := router.Group("/outbound", checkTokenMiddleware)
	routerOutbound.GET("", outboundHandler.Page)
	routerOutbound.POST("", outboundHandler.Create)
	routerOutbound.GET("/:id", outboundHandler.GetById)
	routerOutbound.PUT("/:id", outboundHandler.Update)
	routerOutbound.GET("/:id/set-sent", outboundHandler.SetSent)
	routerOutbound.GET("/:id/generate-delivery-order", outboundHandler.GenerateDeliveryOrder)

	routerStockin := router.Group("/stockin", checkTokenMiddleware)
	routerStockin.GET("", stockinHandler.Page)
	routerStockin.GET("/:id", stockinHandler.GetById)
	routerStockin.POST("", stockinHandler.Create)

	return router

}

func httpErrorHandler(err error, c echo.Context) {
	var errorResponse *response.Response
	code := http.StatusInternalServerError
	switch e := err.(type) {
	case *echo.HTTPError:
		// Handle pada saat URL yang di request tidak ada. atau ada kesalahan server.
		code = e.Code
		errorResponse = &response.Response{
			Status:  false,
			Message: fmt.Sprintf("%v", e.Message),
			Code:    code,
		}
	case *response.Response:
		errorResponse = e
	default:
		// Handle error dari panic
		code = http.StatusInternalServerError
		if config.Debug {
			errorResponse = &response.Response{
				Status:  false,
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		} else {
			errorResponse = &response.Response{
				Status:  false,
				Message: response.ErrorInternalServer,
				Code:    http.StatusInternalServerError,
			}
		}
	}

	js, err := json.Marshal(errorResponse)
	if err == nil {
		_ = c.Blob(code, echo.MIMEApplicationJSON, js)
	} else {
		b := []byte("{status: false, code: 500, message: \"unresolved error\"}")
		_ = c.Blob(code, echo.MIMEApplicationJSON, b)
	}
}

func checkTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		userLogin, err := jwt.ExtractClaims(c.Request().Header.Get(constant.AuthHeaderKey))
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, err.Error()).SendJSON(c)
		}

		conn, closeConn := db.GetConnection()
		defer closeConn()

		var tUser model.User
		err = conn.Where("id = ? ", userLogin.UserID).First(&tUser).Error
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewareUserNotFound).SendJSON(c)
		}

		if tUser.PassVersion != userLogin.PassVersion {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewarePassVersion).SendJSON(c)
		}

		c.Set(constant.TokenUserContext, userLogin)
		return next(c)
	}
}

func checkTokenMiddlewareAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		userLogin, err := jwt.ExtractClaims(c.Request().Header.Get(constant.AuthHeaderKey))
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, err.Error()).SendJSON(c)
		}

		conn, closeConn := db.GetConnection()
		defer closeConn()

		var tUser model.User
		err = conn.Where("id = ? ", userLogin.UserID).First(&tUser).Error
		if err != nil {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewareUserNotFound).SendJSON(c)
		}

		if tUser.PassVersion != userLogin.PassVersion {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorMiddlewarePassVersion).SendJSON(c)
		}

		if tUser.Role == model.UserRoleOperator {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorRoleNotAllowed).SendJSON(c)
		}

		c.Set(constant.TokenUserContext, userLogin)
		return next(c)
	}
}
