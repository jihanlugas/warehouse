package router

import (
	"encoding/json"
	"fmt"

	"github.com/jihanlugas/warehouse/app/auditlog"
	"github.com/jihanlugas/warehouse/app/auth"
	"github.com/jihanlugas/warehouse/app/customer"
	"github.com/jihanlugas/warehouse/app/location"
	"github.com/jihanlugas/warehouse/app/photo"
	"github.com/jihanlugas/warehouse/app/product"
	"github.com/jihanlugas/warehouse/app/purchaseorder"
	"github.com/jihanlugas/warehouse/app/purchaseorderproduct"
	"github.com/jihanlugas/warehouse/app/retail"
	"github.com/jihanlugas/warehouse/app/retailproduct"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stockin"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicle"
	"github.com/jihanlugas/warehouse/app/stockmovementvehiclephoto"
	"github.com/jihanlugas/warehouse/app/stockmovementvehiclepurchaseorder"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicleretail"
	"github.com/jihanlugas/warehouse/app/transaction"
	"github.com/jihanlugas/warehouse/app/transferin"
	"github.com/jihanlugas/warehouse/app/transferout"
	"github.com/jihanlugas/warehouse/app/user"
	"github.com/jihanlugas/warehouse/app/userprivilege"
	"github.com/jihanlugas/warehouse/app/userprovider"
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

	"net/http"

	_ "github.com/jihanlugas/warehouse/docs"
)

func Init() *echo.Echo {
	//// repositories
	auditlogRepository := auditlog.NewRepository()
	userRepository := user.NewRepository()
	userproviderRepository := userprovider.NewRepository()
	userprivilegeRepository := userprivilege.NewRepository()
	locationRepository := location.NewRepository()
	warehouseRepository := warehouse.NewRepository()
	transactionRepository := transaction.NewRepository()
	vehicleRepository := vehicle.NewRepository()
	customerRepository := customer.NewRepository()
	productRepository := product.NewRepository()
	retailRepository := retail.NewRepository()
	retailproductRepository := retailproduct.NewRepository()
	purchaseorderRepository := purchaseorder.NewRepository()
	purchaseorderproductRepository := purchaseorderproduct.NewRepository()
	stockRepository := stock.NewRepository()
	stocklogRepository := stocklog.NewRepository()
	stockmovementvehicleRepository := stockmovementvehicle.NewRepository()
	stockmovementvehiclephotoRepository := stockmovementvehiclephoto.NewRepository()
	photoRepository := photo.NewRepository()

	// usecases
	auditlogUsecase := auditlog.NewUsecase(auditlogRepository)
	authUsecase := auth.NewUsecase(userRepository, warehouseRepository, userproviderRepository)
	userUsecase := user.NewUsecase(userRepository, userprivilegeRepository, warehouseRepository)
	locationUsecase := location.NewUsecase(locationRepository)
	warehouseUsecase := warehouse.NewUsecase(warehouseRepository)
	stockUsecase := stock.NewUsecase(stockRepository, stocklogRepository)
	transactionUsecase := transaction.NewUsecase(transactionRepository)
	vehicleUsecase := vehicle.NewUsecase(vehicleRepository)
	customerUsecase := customer.NewUsecase(customerRepository)
	productUsecase := product.NewUsecase(productRepository)
	stockmovementvehicleUsecase := stockmovementvehicle.NewUsecase(stockmovementvehicleRepository, stockmovementvehiclephotoRepository, photoRepository)
	stockinUsecase := stockin.NewUsecase(stockmovementvehicleRepository, warehouseRepository, stockRepository, stocklogRepository)
	transferoutUsecase := transferout.NewUsecase(stockmovementvehicleRepository, warehouseRepository, stockRepository, stocklogRepository, vehicleRepository)
	transferinUsecase := transferin.NewUsecase(stockmovementvehicleRepository, warehouseRepository, stockRepository, stocklogRepository)
	stockmovementvehiclepurchaseorderUsecase := stockmovementvehiclepurchaseorder.NewUsecase(purchaseorderRepository, stockmovementvehicleRepository, warehouseRepository, stockRepository, stocklogRepository, vehicleRepository)
	stockmovementvehicleretailUsecase := stockmovementvehicleretail.NewUsecase(retailRepository, stockmovementvehicleRepository, warehouseRepository, stockRepository, stocklogRepository, vehicleRepository)
	stocklogUsecase := stocklog.NewUsecase(stocklogRepository)
	//stockmovementvehicleUsecase := stockmovementvehicle.NewUsecase(stockmovementvehicleRepository, warehouseRepository, vehicleRepository, stockRepository, stocklogRepository, purchaseorderRepository, retailRepository, stockmovementvehiclephotoRepository, photoRepository)
	retailUsecase := retail.NewUsecase(retailRepository, retailproductRepository, stockRepository, stocklogRepository, customerRepository, warehouseRepository, vehicleRepository)
	purchaseorderUsecase := purchaseorder.NewUsecase(purchaseorderRepository, purchaseorderproductRepository, stockRepository, stocklogRepository, customerRepository, warehouseRepository, vehicleRepository)

	// handlers
	authHandler := auth.NewHandler(authUsecase, auditlogUsecase)
	auditlogHandler := auditlog.NewHandler(auditlogUsecase)
	userHandler := user.NewHandler(userUsecase)
	locationHandler := location.NewHandler(locationUsecase)
	warehouseHandler := warehouse.NewHandler(warehouseUsecase)
	stockHandler := stock.NewHandler(stockUsecase)
	transactionHandler := transaction.NewHandler(transactionUsecase)
	stocklogHandler := stocklog.NewHandler(stocklogUsecase)
	stockmovementvehicleHandler := stockmovementvehicle.NewHandler(stockmovementvehicleUsecase)
	stockinHandler := stockin.NewHandler(stockinUsecase, auditlogUsecase)
	transferoutHandler := transferout.NewHandler(transferoutUsecase, auditlogUsecase)
	transferinHandler := transferin.NewHandler(transferinUsecase, auditlogUsecase)
	stockmovementvehicleretailHandler := stockmovementvehicleretail.NewHandler(stockmovementvehicleretailUsecase, auditlogUsecase)
	stockmovementvehiclepurchaseorderHandler := stockmovementvehiclepurchaseorder.NewHandler(stockmovementvehiclepurchaseorderUsecase, auditlogUsecase)
	vehicleHandler := vehicle.NewHandler(vehicleUsecase)
	customerHandler := customer.NewHandler(customerUsecase)
	productHandler := product.NewHandler(productUsecase)
	retailHandler := retail.NewHandler(retailUsecase)
	purchaseorderHandler := purchaseorder.NewHandler(purchaseorderUsecase)

	router := websiteRouter()

	router.Static("/"+config.StorageDirectory, config.StorageDirectory)

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
	routerAuth.GET("/google/login", authHandler.GoogleSignIn)
	routerAuth.GET("/google/link", authHandler.GoogleLink, checkTokenMiddlewareQuery)
	routerAuth.GET("/google/callback", authHandler.GoogleCallback)
	routerAuth.GET("/google/unlink", authHandler.GoogleUnlink, checkTokenMiddleware)

	routerAuditlog := router.Group("/auditlog", checkTokenMiddleware)
	routerAuditlog.GET("", auditlogHandler.Page)

	routerLocation := router.Group("/location", checkTokenMiddlewareAdmin)
	routerLocation.GET("", locationHandler.Page)
	routerLocation.GET("/:id", locationHandler.GetById)

	routerCustomer := router.Group("/customer", checkTokenMiddlewareAdmin)
	routerCustomer.GET("", customerHandler.Page)
	routerCustomer.POST("", customerHandler.Create)
	routerCustomer.GET("/:id", customerHandler.GetById)
	routerCustomer.PUT("/:id", customerHandler.Update)
	routerCustomer.DELETE("/:id", customerHandler.Delete)

	routerVehicle := router.Group("/vehicle", checkTokenMiddleware)
	routerVehicle.GET("", vehicleHandler.Page)
	routerVehicle.POST("", vehicleHandler.Create)
	routerVehicle.GET("/:id", vehicleHandler.GetById)
	routerVehicle.PUT("/:id", vehicleHandler.Update)
	routerVehicle.DELETE("/:id", vehicleHandler.Delete)

	routerProduct := router.Group("/product")
	routerProduct.GET("", productHandler.Page, checkTokenMiddleware)
	routerProduct.POST("", productHandler.Create, checkTokenMiddlewareAdmin)
	routerProduct.GET("/:id", productHandler.GetById, checkTokenMiddleware)
	routerProduct.PUT("/:id", productHandler.Update, checkTokenMiddlewareAdmin)
	routerProduct.DELETE("/:id", productHandler.Delete, checkTokenMiddlewareAdmin)

	routerWarehouse := router.Group("/warehouse")
	routerWarehouse.GET("", warehouseHandler.Page, checkTokenMiddleware)
	routerWarehouse.POST("", warehouseHandler.Create, checkTokenMiddlewareAdmin)
	routerWarehouse.GET("/:id", warehouseHandler.GetById, checkTokenMiddleware)
	routerWarehouse.PUT("/:id", warehouseHandler.Update, checkTokenMiddlewareAdmin)
	routerWarehouse.DELETE("/:id", warehouseHandler.Delete, checkTokenMiddlewareAdmin)

	routerStock := router.Group("/stock")
	routerStock.GET("/:id", stockHandler.GetById, checkTokenMiddleware)
	routerStock.PUT("/:id", stockHandler.Update, checkTokenMiddlewareAdmin)

	routerTransaction := router.Group("/transaction", checkTokenMiddlewareAdmin)
	routerTransaction.GET("", transactionHandler.Page)
	routerTransaction.POST("", transactionHandler.Create)
	routerTransaction.GET("/:id", transactionHandler.GetById)
	routerTransaction.PUT("/:id", transactionHandler.Update)
	routerTransaction.DELETE("/:id", transactionHandler.Delete)

	routerUser := router.Group("/user")
	routerUser.GET("", userHandler.Page, checkTokenMiddlewareAdmin)
	routerUser.POST("", userHandler.Create, checkTokenMiddlewareAdmin)
	routerUser.GET("/:id", userHandler.GetById, checkTokenMiddlewareAdmin)
	routerUser.PUT("/:id", userHandler.Update, checkTokenMiddlewareAdmin)
	routerUser.POST("/:id/privilege", userHandler.UpdateUserPrivilege, checkTokenMiddlewareAdmin)
	routerUser.DELETE("/:id", userHandler.Delete, checkTokenMiddlewareAdmin)
	routerUser.POST("/change-password", userHandler.ChangePassword, checkTokenMiddleware)

	routerStockmovementvehicle := router.Group("/stockmovementvehicle")
	routerStockmovementvehicle.GET("", stockmovementvehicleHandler.Page, checkTokenMiddlewareAdmin)
	routerStockmovementvehicle.GET("/:id", stockmovementvehicleHandler.GetById, checkTokenMiddleware)
	routerStockmovementvehicle.DELETE("/:id", stockmovementvehicleHandler.Delete, checkTokenMiddlewareAdmin)
	routerStockmovementvehicle.POST("/:id/upload-photo", stockmovementvehicleHandler.UploadPhoto, checkTokenMiddleware)
	//routerStockmovementvehiclePUTGET("/:id/set-sent", stockmovementvehicleHandler.SetSent, checkTokenMiddleware)
	//routerStockmovementvehicle.GET("/:id/generate-delivery-order", stockmovementvehicleHandler.GenerateDeliveryOrder, checkTokenMiddleware)

	routerStockmovementvehicleStockin := routerStockmovementvehicle.Group("/stock-in", checkTokenMiddleware)
	routerStockmovementvehicleStockin.GET("", stockinHandler.Page)
	routerStockmovementvehicleStockin.POST("", stockinHandler.Create)
	routerStockmovementvehicleStockin.GET("/:id", stockinHandler.GetById)
	routerStockmovementvehicleStockin.DELETE("/:id", stockinHandler.Delete)
	routerStockmovementvehicleStockin.PUT("/:id/set-complete", stockinHandler.SetComplete)

	routerStockmovementvehicleTrasnferout := routerStockmovementvehicle.Group("/transfer-out", checkTokenMiddleware)
	routerStockmovementvehicleTrasnferout.GET("", transferoutHandler.Page)
	routerStockmovementvehicleTrasnferout.POST("", transferoutHandler.Create)
	routerStockmovementvehicleTrasnferout.GET("/:id", transferoutHandler.GetById)
	routerStockmovementvehicleTrasnferout.PUT("/:id", transferoutHandler.Update)
	routerStockmovementvehicleTrasnferout.DELETE("/:id", transferoutHandler.Delete)
	routerStockmovementvehicleTrasnferout.PUT("/:id/set-in-transit", transferoutHandler.SetInTransit)
	routerStockmovementvehicleTrasnferout.GET("/:id/generate-delivery-order", transferoutHandler.GenerateDeliveryOrder)

	routerStockmovementvehicleTrasnferin := routerStockmovementvehicle.Group("/transfer-in", checkTokenMiddleware)
	routerStockmovementvehicleTrasnferin.GET("", transferinHandler.Page)
	routerStockmovementvehicleTrasnferin.GET("/:id", transferinHandler.GetById)
	routerStockmovementvehicleTrasnferin.PUT("/:id", transferinHandler.Update)
	routerStockmovementvehicleTrasnferin.PUT("/:id/set-unloading", transferinHandler.SetUnloading)
	routerStockmovementvehicleTrasnferin.PUT("/:id/set-complete", transferinHandler.SetComplete)
	routerStockmovementvehicleTrasnferin.GET("/:id/generate-delivery-recipt", transferinHandler.GenerateDeliveryRecipt)

	routerStockmovementvehiclePurchaseorder := routerStockmovementvehicle.Group("/purchase-order", checkTokenMiddleware)
	routerStockmovementvehiclePurchaseorder.GET("", stockmovementvehiclepurchaseorderHandler.Page)
	routerStockmovementvehiclePurchaseorder.POST("", stockmovementvehiclepurchaseorderHandler.Create)
	routerStockmovementvehiclePurchaseorder.GET("/:id", stockmovementvehiclepurchaseorderHandler.GetById)
	routerStockmovementvehiclePurchaseorder.PUT("/:id", stockmovementvehiclepurchaseorderHandler.Update)
	routerStockmovementvehiclePurchaseorder.DELETE("/:id", stockmovementvehiclepurchaseorderHandler.Delete)
	routerStockmovementvehiclePurchaseorder.PUT("/:id/set-complete", stockmovementvehiclepurchaseorderHandler.SetComplete)
	routerStockmovementvehiclePurchaseorder.GET("/:id/generate-delivery-order", stockmovementvehiclepurchaseorderHandler.GenerateDeliveryOrder)

	routerStockmovementvehicleRetail := routerStockmovementvehicle.Group("/retail", checkTokenMiddleware)
	routerStockmovementvehicleRetail.GET("", stockmovementvehicleretailHandler.Page)
	routerStockmovementvehicleRetail.POST("", stockmovementvehicleretailHandler.Create)
	routerStockmovementvehicleRetail.GET("/:id", stockmovementvehicleretailHandler.GetById)
	routerStockmovementvehicleRetail.PUT("/:id", stockmovementvehicleretailHandler.Update)
	routerStockmovementvehicleRetail.DELETE("/:id", stockmovementvehicleretailHandler.Delete)
	routerStockmovementvehicleRetail.PUT("/:id/set-complete", stockmovementvehicleretailHandler.SetComplete)
	routerStockmovementvehicleRetail.GET("/:id/generate-delivery-order", stockmovementvehicleretailHandler.GenerateDeliveryOrder)

	routerStocklog := router.Group("/stocklog")
	routerStocklog.GET("", stocklogHandler.Page, checkTokenMiddleware)
	routerStocklog.GET("/:id", stocklogHandler.GetById, checkTokenMiddleware)

	//
	//routerStockmovementvehiclePurchaseorder := routerStockmovementvehicle.Group("/purchaseorder")
	//routerStockmovementvehiclePurchaseorder.POST("", stockmovementvehicleHandler.CreatePurchaseorder, checkTokenMiddleware)
	//routerStockmovementvehiclePurchaseorder.PUT("/:id", stockmovementvehicleHandler.UpdatePurchaseorder, checkTokenMiddleware)
	//
	//routerStockmovementvehicleRetail := routerStockmovementvehicle.Group("/retail")
	//routerStockmovementvehicleRetail.POST("", stockmovementvehicleHandler.CreateRetail, checkTokenMiddleware)
	//routerStockmovementvehicleRetail.PUT("/:id", stockmovementvehicleHandler.UpdateRetail, checkTokenMiddleware)

	routerRetail := router.Group("/retail", checkTokenMiddleware)
	routerRetail.GET("", retailHandler.Page)
	routerRetail.POST("", retailHandler.Create)
	routerRetail.GET("/:id", retailHandler.GetById)
	routerRetail.PUT("/:id", retailHandler.Update)
	routerRetail.DELETE("/:id", retailHandler.Delete)
	routerRetail.PUT("/:id/set-status-open", retailHandler.SetStatusOpen)
	routerRetail.PUT("/:id/set-status-close", retailHandler.SetStatusClose)
	routerRetail.GET("/:id/generate-invoice", retailHandler.GenerateInvoice)

	routerPurchaseorder := router.Group("/purchase-order", checkTokenMiddleware)
	routerPurchaseorder.GET("", purchaseorderHandler.Page)
	routerPurchaseorder.POST("", purchaseorderHandler.Create)
	routerPurchaseorder.GET("/:id", purchaseorderHandler.GetById)
	routerPurchaseorder.PUT("/:id", purchaseorderHandler.Update)
	routerPurchaseorder.DELETE("/:id", purchaseorderHandler.Delete)
	routerPurchaseorder.PUT("/:id/set-status-open", purchaseorderHandler.SetStatusOpen)
	routerPurchaseorder.PUT("/:id/set-status-close", purchaseorderHandler.SetStatusClose)
	routerPurchaseorder.GET("/:id/generate-invoice", purchaseorderHandler.GenerateInvoice)

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

func checkTokenMiddlewareQuery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		token := c.QueryParam("token")

		userLogin, err := jwt.ExtractClaimsQuery(token)
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

		if tUser.UserRole == model.UserRoleOperator {
			return response.ErrorForce(http.StatusUnauthorized, response.ErrorRoleNotAllowed).SendJSON(c)
		}

		c.Set(constant.TokenUserContext, userLogin)
		return next(c)
	}
}
