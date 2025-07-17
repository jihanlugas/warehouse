package stockmovementvehicle

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/purchaseorder"
	"github.com/jihanlugas/warehouse/app/retail"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovement"
	"github.com/jihanlugas/warehouse/app/vehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"html/template"
	"os"
	"strings"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Delete(loginUser jwt.UserLogin, id string) error
	SetSent(loginUser jwt.UserLogin, id string) error
	GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error)
	CreatePurchaseorder(loginUser jwt.UserLogin, req request.CreateStockmovementvehiclePurchaseorder) error
	UpdatePurchaseorder(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehiclePurchaseorder) error
	CreateRetail(loginUser jwt.UserLogin, req request.CreateStockmovementvehicleRetail) error
	UpdateRetail(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehicleRetail) error
}

type usecase struct {
	stockmovementvehicleRepository Repository
	warehouseRepository            warehouse.Repository
	vehicleRepository              vehicle.Repository
	stockmovementRepository        stockmovement.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	purchaseorderRepository        purchaseorder.Repository
	retailRepository               retail.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicles, count, err = u.stockmovementvehicleRepository.Page(conn, req)
	if err != nil {
		return vStockmovementvehicles, count, err
	}

	return vStockmovementvehicles, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	return vStockmovementvehicle, err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if tStockmovementvehicle.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	tx := conn.Begin()

	err = u.stockmovementvehicleRepository.Delete(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetSent(loginUser jwt.UserLogin, id string) error {
	var err error
	var vStockmovementvehicle model.StockmovementvehicleView
	var tStock model.Stock
	var tStocklog model.Stocklog
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if tStockmovementvehicle.SentNetQuantity <= 0 && tStockmovementvehicle.SentGrossQuantity <= 0 {
		return errors.New("data gross or net quantity is zero")
	}

	tx := conn.Begin()

	now := time.Now()
	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, vStockmovementvehicle.FromWarehouseID, vStockmovementvehicle.ProductID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
		}
		tStock = model.Stock{
			ID:          utils.GetUniqueID(),
			WarehouseID: vStockmovementvehicle.FromWarehouseID,
			ProductID:   vStockmovementvehicle.ProductID,
			Quantity:    0,
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
		}
	}

	tStockmovementvehicle.SentTime = &now
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	CurrentQuantity := 0.0
	if vStockmovementvehicle.SentNetQuantity != 0 {
		CurrentQuantity = tStock.Quantity - vStockmovementvehicle.SentNetQuantity
		tStock.Quantity = CurrentQuantity
		tStock.UpdateBy = loginUser.UserID
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockRepository.Name(), err))
		}

		tStocklog = model.Stocklog{
			WarehouseID:            vStockmovementvehicle.FromWarehouseID,
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
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stocklogRepository.Name(), err))
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, "Stockmovement", "Vehicle", "Product", "Purchaseorder", "Purchaseorder.Customer", "Retail", "Retail.Customer")
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.Stockmovement == nil {
		log.Info("stockmovement not found")
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorDataNotFound)
	}

	if vStockmovementvehicle.FromWarehouseID == "" {
		log.Info("warehouse not found")
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorDataNotFound)
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	switch vStockmovementvehicle.Stockmovement.Type {
	case model.StockMovementTypePurchaseOrder:
		pdfBytes, err = u.generatePurchaseorderDeliveryOrder(vStockmovementvehicle)
	case model.StockMovementTypeRetail:
		pdfBytes, err = u.generateRetailDeliveryOrder(vStockmovementvehicle)
	case model.StockMovementTypeTransfer:
		pdfBytes, err = u.generateTransferDeliveryOrder(vStockmovementvehicle)
	default:
		return pdfBytes, vStockmovementvehicle, errors.New("stockmovement type not recognized")
	}

	return pdfBytes, vStockmovementvehicle, err
}

func (u usecase) generatePurchaseorderDeliveryOrder(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("purchaseorder-delivery-order.html").Funcs(template.FuncMap{
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
		"displaySpkNumber":    utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/purchaseorder-delivery-order.html")
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

func (u usecase) generateRetailDeliveryOrder(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("retail-delivery-order.html").Funcs(template.FuncMap{
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
		"displaySpkNumber":    utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/retail-delivery-order.html")
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

func (u usecase) generateTransferDeliveryOrder(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("delivery-order.html").Funcs(template.FuncMap{
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

func (u usecase) CreatePurchaseorder(loginUser jwt.UserLogin, req request.CreateStockmovementvehiclePurchaseorder) error {
	var err error
	var vPurchaseorder model.PurchaseorderView
	var vWarehouse model.WarehouseView
	var tVehicle model.Vehicle
	var vStockmovement model.StockmovementView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPurchaseorder, err = u.purchaseorderRepository.GetViewById(conn, req.PurchaseorderID, "Purchaseorderproducts")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	if vPurchaseorder.Status != model.PurchaseorderStatusOpen {
		return errors.New("purchase order is not open")
	}

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, req.FromWarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsPurchaseorder {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create purchaseorder"))
	}

	tx := conn.Begin()

	vStockmovement, err = u.stockmovementRepository.GetViewByFromWarehouseIDAndRelatedIDAndProductID(tx, vWarehouse.ID, req.PurchaseorderID, req.ProductID)
	if err != nil {
		var vPurchaseorderproduct model.PurchaseorderproductView
		found := false
		for _, purchaseorderproduct := range vPurchaseorder.Purchaseorderproducts {
			if purchaseorderproduct.ProductID == req.ProductID {
				vPurchaseorderproduct = purchaseorderproduct
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("failed to get purchaseorder product %s: %v", req.ProductID, err))
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			tStockmovement := model.Stockmovement{
				ID:              utils.GetUniqueID(),
				FromWarehouseID: vWarehouse.ID,
				ProductID:       req.ProductID,
				RelatedID:       req.PurchaseorderID,
				Type:            model.StockMovementTypePurchaseOrder,
				UnitPrice:       vPurchaseorderproduct.UnitPrice,
				Remark:          "",
				CreateBy:        "",
				CreateDt:        time.Time{},
				UpdateBy:        "",
				UpdateDt:        time.Time{},
				DeleteDt:        gorm.DeletedAt{},
			}
			err = u.stockmovementRepository.Create(tx, tStockmovement)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementRepository.Name(), err))
			}

			vStockmovement, err = u.stockmovementRepository.GetViewByFromWarehouseIDAndRelatedIDAndProductID(tx, vWarehouse.ID, req.PurchaseorderID, req.ProductID)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementRepository.Name(), err))
			}
		} else {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementRepository.Name(), err))
		}
	}

	// todo: check purchaseorder still open

	if req.IsNewVehiclerdriver {
		tVehicle = model.Vehicle{
			ID:          utils.GetUniqueID(),
			WarehouseID: req.FromWarehouseID,
			PlateNumber: strings.ToUpper(req.PlateNumber),
			Name:        req.VehicleName,
			NIK:         req.NIK,
			DriverName:  req.DriverName,
			PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.vehicleRepository.Create(tx, tVehicle)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.vehicleRepository.Name(), err))
		}
	} else {
		tVehicle, err = u.vehicleRepository.GetTableById(conn, req.VehicleID)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
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

func (u usecase) UpdatePurchaseorder(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehiclePurchaseorder) error {
	var err error
	var vStockmovementvehicle model.StockmovementvehicleView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tStockmovementvehicle.SentGrossQuantity = req.SentGrossQuantity
	tStockmovementvehicle.SentTareQuantity = req.SentTareQuantity
	tStockmovementvehicle.SentNetQuantity = req.SentNetQuantity
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) CreateRetail(loginUser jwt.UserLogin, req request.CreateStockmovementvehicleRetail) error {
	var err error
	var vRetail model.RetailView
	var vWarehouse model.WarehouseView
	var tVehicle model.Vehicle
	var vStockmovement model.StockmovementView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vRetail, err = u.retailRepository.GetViewById(conn, req.RetailID, "Retailproducts")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.retailRepository.Name(), err))
	}

	if vRetail.Status != model.RetailStatusOpen {
		return errors.New("retail is not open")
	}

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, req.FromWarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsRetail {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create retail"))
	}

	tx := conn.Begin()

	vStockmovement, err = u.stockmovementRepository.GetViewByFromWarehouseIDAndRelatedIDAndProductID(tx, vWarehouse.ID, req.RetailID, req.ProductID)
	if err != nil {
		var vRetailproduct model.RetailproductView
		found := false
		for _, retailproduct := range vRetail.Retailproducts {
			if retailproduct.ProductID == req.ProductID {
				vRetailproduct = retailproduct
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("failed to get retail product %s: %v", req.ProductID, err))
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			tStockmovement := model.Stockmovement{
				ID:              utils.GetUniqueID(),
				FromWarehouseID: vWarehouse.ID,
				ProductID:       req.ProductID,
				RelatedID:       req.RetailID,
				Type:            model.StockMovementTypeRetail,
				UnitPrice:       vRetailproduct.UnitPrice,
				Remark:          "",
				CreateBy:        "",
				CreateDt:        time.Time{},
				UpdateBy:        "",
				UpdateDt:        time.Time{},
				DeleteDt:        gorm.DeletedAt{},
			}
			err = u.stockmovementRepository.Create(tx, tStockmovement)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementRepository.Name(), err))
			}

			vStockmovement, err = u.stockmovementRepository.GetViewByFromWarehouseIDAndRelatedIDAndProductID(tx, vWarehouse.ID, req.RetailID, req.ProductID)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementRepository.Name(), err))
			}
		} else {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementRepository.Name(), err))
		}
	}

	// todo: check retail still open

	if req.IsNewVehiclerdriver {
		tVehicle = model.Vehicle{
			ID:          utils.GetUniqueID(),
			WarehouseID: req.FromWarehouseID,
			PlateNumber: strings.ToUpper(req.PlateNumber),
			Name:        req.VehicleName,
			NIK:         req.NIK,
			DriverName:  req.DriverName,
			PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.vehicleRepository.Create(tx, tVehicle)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.vehicleRepository.Name(), err))
		}
	} else {
		tVehicle, err = u.vehicleRepository.GetTableById(conn, req.VehicleID)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
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

func (u usecase) UpdateRetail(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehicleRetail) error {
	var err error
	var vStockmovementvehicle model.StockmovementvehicleView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tStockmovementvehicle.SentGrossQuantity = req.SentGrossQuantity
	tStockmovementvehicle.SentTareQuantity = req.SentTareQuantity
	tStockmovementvehicle.SentNetQuantity = req.SentNetQuantity
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(stockmovementvehicleRepository Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository, stockmovementRepository stockmovement.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, purchaseorderRepository purchaseorder.Repository, retailRepository retail.Repository) Usecase {
	return &usecase{
		stockmovementvehicleRepository: stockmovementvehicleRepository,
		warehouseRepository:            warehouseRepository,
		vehicleRepository:              vehicleRepository,
		stockmovementRepository:        stockmovementRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		purchaseorderRepository:        purchaseorderRepository,
		retailRepository:               retailRepository,
	}
}
