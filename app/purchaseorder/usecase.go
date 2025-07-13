package purchaseorder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/customer"
	"github.com/jihanlugas/warehouse/app/purchaseorderproduct"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/vehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
	"html/template"
	"os"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PagePurchaseorder) (vPurchaseorders []model.PurchaseorderView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vPurchaseorder model.PurchaseorderView, err error)
	Create(loginUser jwt.UserLogin, req request.CreatePurchaseorder) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdatePurchaseorder) error
	Delete(loginUser jwt.UserLogin, id string) error
	SetStatusOpen(loginUser jwt.UserLogin, id string) error
	SetStatusClose(loginUser jwt.UserLogin, id string) error
	GenerateInvoice(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vPurchaseorder model.PurchaseorderView, err error)
}

type usecase struct {
	purchaseorderRepository        Repository
	purchaseorderproductRepository purchaseorderproduct.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	customerRepository             customer.Repository
	warehouseRepository            warehouse.Repository
	vehicleRepository              vehicle.Repository
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
		return vPurchaseorder, errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	return vPurchaseorder, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreatePurchaseorder) error {
	var err error
	var tCustomer model.Customer
	var tPurchaseorder model.Purchaseorder
	var tPurchaseorderproduct model.Purchaseorderproduct

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
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.customerRepository.Name(), err))
		}
		req.CustomerID = tCustomer.ID
	}

	tPurchaseorder = model.Purchaseorder{
		ID:         utils.GetUniqueID(),
		CustomerID: req.CustomerID,
		Notes:      req.Notes,
		Status:     model.PurchaseorderStatusOpen,
		CreateBy:   loginUser.UserID,
		UpdateBy:   loginUser.UserID,
	}

	err = u.purchaseorderRepository.Create(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.purchaseorderRepository.Name(), err))
	}

	for _, product := range req.Products {
		tPurchaseorderproduct = model.Purchaseorderproduct{
			ID:              utils.GetUniqueID(),
			PurchaseorderID: tPurchaseorder.ID,
			ProductID:       product.ProductID,
			UnitPrice:       product.UnitPrice,
			CreateBy:        loginUser.UserID,
			UpdateBy:        loginUser.UserID,
		}
		err = u.purchaseorderproductRepository.Create(tx, tPurchaseorderproduct)
		if err != nil {
			return fmt.Errorf("failed to create %s: %v", u.purchaseorderproductRepository.Name(), err)
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
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	tx := conn.Begin()

	tPurchaseorder.Notes = req.Notes
	tPurchaseorder.UpdateBy = loginUser.UserID
	err = u.purchaseorderRepository.Save(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.purchaseorderRepository.Name(), err))
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
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	tx := conn.Begin()

	err = u.purchaseorderRepository.Delete(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.purchaseorderRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetStatusOpen(loginUser jwt.UserLogin, id string) error {
	var err error
	var tPurchaseorder model.Purchaseorder

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPurchaseorder, err = u.purchaseorderRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	tx := conn.Begin()

	if tPurchaseorder.Status == model.PurchaseorderStatusOpen {
		return errors.New("the purchase order is already open")
	}

	tPurchaseorder.Status = model.PurchaseorderStatusOpen
	tPurchaseorder.UpdateBy = loginUser.UserID
	err = u.purchaseorderRepository.Save(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set status open %s: %v", u.purchaseorderRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetStatusClose(loginUser jwt.UserLogin, id string) error {
	var err error
	var tPurchaseorder model.Purchaseorder

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tPurchaseorder, err = u.purchaseorderRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	tx := conn.Begin()

	if tPurchaseorder.Status == model.PurchaseorderStatusClose {
		return errors.New("the purchase order is already open")
	}

	tPurchaseorder.Status = model.PurchaseorderStatusClose
	tPurchaseorder.UpdateBy = loginUser.UserID
	err = u.purchaseorderRepository.Save(tx, tPurchaseorder)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set status close %s: %v", u.purchaseorderRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateInvoice(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vPurchaseorder model.PurchaseorderView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vPurchaseorder, err = u.purchaseorderRepository.GetViewById(conn, id, "Customer", "Stockmovementvehicles", "Stockmovementvehicles.Stockmovement", "Stockmovementvehicles.Product", "Transactions")
	if err != nil {
		return pdfBytes, vPurchaseorder, errors.New(fmt.Sprintf("failed to get %s: %v", u.purchaseorderRepository.Name(), err))
	}

	pdfBytes, err = u.generateInvoice(vPurchaseorder)

	return pdfBytes, vPurchaseorder, err
}

func (u usecase) generateInvoice(vPurchaseorder model.PurchaseorderView) (pdfBytes []byte, err error) {
	tmpl := template.New("purchaseorder-invoice.html").Funcs(template.FuncMap{
		"displayNumberMinus": func(a, b float64) string {
			return utils.DisplayNumber(a - b)
		},
		"displayMoneyMultiple": func(a, b float64) string {
			return utils.DisplayMoney(a * b)
		},
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/purchaseorder-invoice.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vPurchaseorder); err != nil {
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

func NewUsecase(purchaseorderRepository Repository, purchaseorderproductRepository purchaseorderproduct.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, customerRepository customer.Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository) Usecase {
	return &usecase{
		purchaseorderRepository:        purchaseorderRepository,
		purchaseorderproductRepository: purchaseorderproductRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		customerRepository:             customerRepository,
		warehouseRepository:            warehouseRepository,
		vehicleRepository:              vehicleRepository,
	}
}
