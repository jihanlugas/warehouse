package purchaseorder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/customer"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovement"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicle"
	"github.com/jihanlugas/warehouse/app/vehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"html/template"
	"os"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PagePurchaseorder) (vPurchaseorders []model.PurchaseorderView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPurchaseorder model.PurchaseorderView, err error)
	Create(loginUser jwt.UserLogin, req request.CreatePurchaseorder) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdatePurchaseorder) error
	Delete(loginUser jwt.UserLogin, id string) error
	DeliveryPage(loginUser jwt.UserLogin, id string, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	DeliveryGetById(loginUser jwt.UserLogin, id string, stockmovementvehicleId string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	DeliveryCreate(loginUser jwt.UserLogin, id string, req request.CreatePurchaseorderStockmovementvehicle) error
	DeliveryUpdate(loginUser jwt.UserLogin, id string, stockmovementvehicleId string, req request.UpdatePurchaseorderStockmomentvehicle) error
	DeliveryGenerateDeliveryOrder(loginUser jwt.UserLogin, id string, stockmovementvehicleId string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error)
}

type usecase struct {
	purchaseorderRepository        Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	customerRepository             customer.Repository
	warehouseRepository            warehouse.Repository
	vehicleRepository              vehicle.Repository
	stockmovementRepository        stockmovement.Repository
	stockmovementvehicleRepository stockmovementvehicle.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PagePurchaseorder) (vPurchaseorders []model.PurchaseorderView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPurchaseorders, count, err = u.purchaseorderRepository.Page(conn, req)
	if err != nil {
		return vPurchaseorders, count, err
	}

	return vPurchaseorders, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPurchaseorder model.PurchaseorderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPurchaseorder, err = u.purchaseorderRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vPurchaseorder, errors.New(fmt.Sprint("failed to get purchaseorder: ", err))
	}

	return vPurchaseorder, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreatePurchaseorder) error {
	var err error
	var tCustomer model.Customer
	var tPurchaseorder model.Purchaseorder
	var tStockmovement model.Stockmovement
	var vWarehouses []model.WarehouseView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	if req.IsNewCustomer {
		tCustomer = model.Customer{
			ID:          utils.GetUniqueID(),
			Name:        req.CustomerName,
			PhoneNumber: utils.FormatPhoneTo62(req.CustomerPhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}

		err = u.customerRepository.Create(tx, tCustomer)
		if err != nil {
			return errors.New(fmt.Sprint("failed to create customer: ", err))
		}
		req.CustomerID = tCustomer.ID
	}

	tPurchaseorder = model.Purchaseorder{
		ID:          utils.GetUniqueID(),
		CustomerID:  req.CustomerID,
		TotalAmount: 0,
		Notes:       req.Notes,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.purchaseorderRepository.Create(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create purchaseorder: ", err))
	}

	vWarehousesReq := request.PageWarehouse{
		Paging: request.Paging{
			Limit: -1,
		},
	}
	vWarehouses, _, err = u.warehouseRepository.Page(tx, vWarehousesReq)
	if err != nil {
		return fmt.Errorf("failed to get page %s: %v", u.warehouseRepository.Name(), err)
	}

	for _, product := range req.Products {
		for _, vWarehouse := range vWarehouses {
			if vWarehouse.IsPurchaseorder {
				tStockmovement = model.Stockmovement{
					ID:              utils.GetUniqueID(),
					FromWarehouseID: vWarehouse.ID,
					ProductID:       product.ProductID,
					RelatedID:       tPurchaseorder.ID,
					Type:            model.StockMovementTypePurchaseOrder,
					UnitPrice:       product.UnitPrice,
					CreateBy:        loginUser.UserID,
					UpdateBy:        loginUser.UserID,
				}
				err = u.stockmovementRepository.Create(tx, tStockmovement)
				if err != nil {
					return fmt.Errorf("failed to create page %s: %v", u.stockmovementRepository.Name(), err)
				}
			}
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdatePurchaseorder) error {
	var err error
	var tPurchaseorder model.Purchaseorder

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPurchaseorder, err = u.purchaseorderRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get purchaseorder: ", err))
	}

	if tPurchaseorder.TotalAmount != 0 {
		return errors.New("unable to update data")
	}

	tx := conn.Begin()

	tPurchaseorder.Notes = req.Notes
	tPurchaseorder.UpdateBy = loginUser.UserID
	err = u.purchaseorderRepository.Save(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update purchaseorder: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tPurchaseorder model.Purchaseorder

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPurchaseorder, err = u.purchaseorderRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get purchaseorder: ", err))
	}

	tx := conn.Begin()

	err = u.purchaseorderRepository.Delete(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete purchaseorder: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) DeliveryPage(loginUser jwt.UserLogin, id string, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicles, count, err = u.stockmovementvehicleRepository.Page(conn, req)
	if err != nil {
		return vStockmovementvehicles, count, err
	}

	return vStockmovementvehicles, count, err
}

func (u usecase) DeliveryGetById(loginUser jwt.UserLogin, id string, stockmovementvehicleId string, preloads ...string) (vStockmomentvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmomentvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, stockmovementvehicleId, preloads...)
	if err != nil {
		return vStockmomentvehicle, errors.New(fmt.Sprint("failed to get delivery purchaseorder: ", err))
	}

	return vStockmomentvehicle, err
}

func (u usecase) DeliveryCreate(loginUser jwt.UserLogin, id string, req request.CreatePurchaseorderStockmovementvehicle) error {
	var err error
	var vPurchaseorder model.PurchaseorderView
	var vWarehouse model.WarehouseView
	var tVehicle model.Vehicle
	var vStockmovement model.StockmovementView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	vPurchaseorder, err = u.purchaseorderRepository.GetViewById(tx, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get purchaseorder: ", err))
	}

	vStockmovement, err = u.stockmovementRepository.GetViewByRelatedIDAndProductID(tx, vPurchaseorder.ID, req.ProductID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get stockmovement: ", err))
	}

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, vStockmovement.FromWarehouseID)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get warehouse: ", err))
	}

	if !vWarehouse.IsPurchaseorder {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create purchaseorder"))
	}

	if req.IsNewVehiclerdriver {
		tVehicle = model.Vehicle{
			ID:          utils.GetUniqueID(),
			WarehouseID: vWarehouse.ID,
			PlateNumber: req.PlateNumber,
			Name:        req.VehicleName,
			NIK:         req.NIK,
			DriverName:  req.DriverName,
			PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.vehicleRepository.Create(tx, tVehicle)
		if err != nil {
			return errors.New(fmt.Sprint("failed to create vehicle: ", err))
		}
	} else {
		tVehicle, err = u.vehicleRepository.GetTableById(conn, req.VehicleID)
		if err != nil {
			return errors.New(fmt.Sprint("failed to create vehicle : ", err))
		}
	}

	tStockmovementvehicle = model.Stockmovementvehicle{
		StockmovementID:   vStockmovement.ID,
		ProductID:         req.ProductID,
		VehicleID:         tVehicle.ID,
		SentGrossQuantity: req.SentGrossQuantity,
		SentTareQuantity:  req.SentTareQuantity,
		SentNetQuantity:   req.SentNetQuantity,
		CreateBy:          loginUser.UserID,
		UpdateBy:          loginUser.UserID,
	}
	err = u.stockmovementvehicleRepository.Create(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create stockmovement vehicle: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) DeliveryUpdate(loginUser jwt.UserLogin, id string, stockmovementvehicleId string, req request.UpdatePurchaseorderStockmomentvehicle) error {
	var err error
	var vStockmovement model.StockmovementView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get stockmovementvehicle: ", err))
	}

	vStockmovement, err = u.stockmovementRepository.GetViewById(conn, tStockmovementvehicle.StockmovementID)

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovement.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	tx := conn.Begin()

	tStockmovementvehicle.SentGrossQuantity = req.SentGrossQuantity
	tStockmovementvehicle.SentTareQuantity = req.SentTareQuantity
	tStockmovementvehicle.SentNetQuantity = req.SentNetQuantity
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update stockmovementvehicle: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err

}

func (u usecase) DeliveryGenerateDeliveryOrder(loginUser jwt.UserLogin, id string, stockmovementvehicleId string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error) {
	var tStockmovementvehicle model.Stockmovementvehicle
	var tStocklog model.Stocklog
	var tStock model.Stock
	var vPurchaseorder model.PurchaseorderView
	var vWarehouse model.WarehouseView

	now := time.Now()
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPurchaseorder, err = u.purchaseorderRepository.GetViewById(conn, id, "Stockmovements")
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to get purchaseorder: ", err))
	}

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, vPurchaseorder.Stockmovements[0].FromWarehouseID)
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to get warehouse: ", err))
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to get stockmovementvehicle: ", err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vWarehouse.ID) {
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.RecivedNetQuantity != 0 || tStockmovementvehicle.RecivedGrossQuantity != 0 {
		return pdfBytes, vStockmovementvehicle, errors.New("data recived or net quantity is zero")
	}

	tx := conn.Begin()

	// set sent time and insert stock log, calculate stock if warehouse isStock true
	if tStockmovementvehicle.SentTime == nil {
		tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, vWarehouse.ID, vStockmovementvehicle.ProductID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get stock %s: %v", u.stockRepository.Name(), err))
			}
			tStock = model.Stock{
				ID:          utils.GetUniqueID(),
				WarehouseID: vWarehouse.ID,
				ProductID:   vStockmovementvehicle.ProductID,
				Quantity:    0,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			err = u.stockRepository.Save(tx, tStock)
			if err != nil {
				return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
			}
		}

		tStockmovementvehicle.SentTime = &now
		tStockmovementvehicle.UpdateBy = loginUser.UserID
		err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
		if err != nil {
			return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to update stockmovementvehicle: ", err))
		}

		CurrentQuantity := 0.0
		if vStockmovementvehicle.SentNetQuantity != 0 {
			CurrentQuantity = tStock.Quantity - vStockmovementvehicle.SentNetQuantity
			tStock.Quantity = CurrentQuantity
			tStock.UpdateBy = loginUser.UserID
			err = u.stockRepository.Save(tx, tStock)
			if err != nil {
				return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to update stock: ", err))
			}

			tStocklog = model.Stocklog{
				WarehouseID:            vWarehouse.ID,
				StockID:                tStock.ID,
				StockmovementID:        tStockmovementvehicle.StockmovementID,
				StockmovementvehicleID: tStockmovementvehicle.ID,
				ProductID:              tStockmovementvehicle.ProductID,
				VehicleID:              tStockmovementvehicle.VehicleID,
				Type:                   model.StockLogTypeOut,
				GrossQuantity:          vStockmovementvehicle.SentGrossQuantity,
				TareQuantity:           vStockmovementvehicle.SentTareQuantity,
				NetQuantity:            vStockmovementvehicle.SentNetQuantity,
				CurrentQuantity:        CurrentQuantity,
				CreateBy:               loginUser.UserID,
				UpdateBy:               loginUser.UserID,
			}
			err = u.stocklogRepository.Create(tx, tStocklog)
			if err != nil {
				return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to create stocklog: ", err))
			}
		}
		vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id)
		if err != nil {
			return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprint("failed to get deliverypurchaseorder: ", err))
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return pdfBytes, vStockmovementvehicle, err
	}

	pdfBytes, err = u.generateDeliveryOrder(vStockmovementvehicle)

	return pdfBytes, vStockmovementvehicle, err
}

func (u usecase) generateDeliveryOrder(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("purchaseorder-delivery-order.html").Funcs(template.FuncMap{
		"displayLembar": func(lembar int64) string {
			return fmt.Sprintf("%s Lembar", utils.DisplayNumber(lembar))
		},
		"displayDuplex": func(isDuplex bool) string {
			if isDuplex {
				return "2 Muka"
			}
			return "1 Muka"
		},
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
		"displaySpkNumber":    utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/delivery-order.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vStockmovementvehicle); err != nil {
		return pdfBytes, err
	}

	// Simpan HTML render ke file sementara
	tempHTMLFile := "temp.html"
	if err := os.WriteFile(tempHTMLFile, buf.Bytes(), 0644); err != nil {
		return pdfBytes, err
	}
	defer os.Remove(tempHTMLFile)

	return utils.GeneratePDFWithChromedp(tempHTMLFile)
}

func NewUsecase(purchaseorderRepository Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, customerRepository customer.Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository, stockmovementRepository stockmovement.Repository, stockmovementvehicleRepository stockmovementvehicle.Repository) Usecase {
	return &usecase{
		purchaseorderRepository:        purchaseorderRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		customerRepository:             customerRepository,
		warehouseRepository:            warehouseRepository,
		vehicleRepository:              vehicleRepository,
		stockmovementRepository:        stockmovementRepository,
		stockmovementvehicleRepository: stockmovementvehicleRepository,
	}
}
