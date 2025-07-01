package cmd

import (
	"github.com/jihanlugas/warehouse/cryption"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

func dbUp() {
	log.Info("Running database migrations...")
	dbUpTable()
	dbUpView()
}

func dbUpTable() {
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Userprivilege{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Customer{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Retail{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Purchaseorder{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Transaction{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Warehouse{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Stock{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Stocklog{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Vehicle{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Product{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Stockmovement{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Stockmovementvehicle{})
	if err != nil {
		panic(err)
	}
}

func dbUpView() {
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Migrator().DropView(model.VIEW_USER)
	if err != nil {
		panic(err)
	}
	vUser := conn.Model(&model.User{}).Unscoped().
		Select("users.*, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = users.create_by").
		Joins("left join users u2 on u2.id = users.update_by")
	err = conn.Migrator().CreateView(model.VIEW_USER, gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_USERPRIVILEGE)
	if err != nil {
		panic(err)
	}
	vUserprivilege := conn.Model(&model.Userprivilege{}).Unscoped().
		Select("userprivileges.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = userprivileges.create_by").
		Joins("left join users u2 on u2.id = userprivileges.update_by")
	err = conn.Migrator().CreateView(model.VIEW_USERPRIVILEGE, gorm.ViewOption{
		Replace: true,
		Query:   vUserprivilege,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_CUSTOMER)
	if err != nil {
		panic(err)
	}
	vCustomer := conn.Model(&model.Customer{}).Unscoped().
		Select("customers.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = customers.create_by").
		Joins("left join users u2 on u2.id = customers.update_by")
	err = conn.Migrator().CreateView(model.VIEW_CUSTOMER, gorm.ViewOption{
		Replace: true,
		Query:   vCustomer,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_RETAIL)
	if err != nil {
		panic(err)
	}
	vRetail := conn.Model(&model.Retail{}).Unscoped().
		Select("retails.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = retails.create_by").
		Joins("left join users u2 on u2.id = retails.update_by")
	err = conn.Migrator().CreateView(model.VIEW_RETAIL, gorm.ViewOption{
		Replace: true,
		Query:   vRetail,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_PURCHASEORDER)
	if err != nil {
		panic(err)
	}
	vPurchaseorder := conn.Model(&model.Purchaseorder{}).Unscoped().
		Select("purchaseorders.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = purchaseorders.create_by").
		Joins("left join users u2 on u2.id = purchaseorders.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PURCHASEORDER, gorm.ViewOption{
		Replace: true,
		Query:   vPurchaseorder,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_TRANSACTION)
	if err != nil {
		panic(err)
	}
	vTransaction := conn.Model(&model.Transaction{}).Unscoped().
		Select("transactions.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = transactions.create_by").
		Joins("left join users u2 on u2.id = transactions.update_by")
	err = conn.Migrator().CreateView(model.VIEW_TRANSACTION, gorm.ViewOption{
		Replace: true,
		Query:   vTransaction,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_WAREHOUSE)
	if err != nil {
		panic(err)
	}
	vWarehouse := conn.Model(&model.Warehouse{}).Unscoped().
		Select("warehouses.*, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = warehouses.create_by").
		Joins("left join users u2 on u2.id = warehouses.update_by")
	err = conn.Migrator().CreateView(model.VIEW_WAREHOUSE, gorm.ViewOption{
		Replace: true,
		Query:   vWarehouse,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCK)
	if err != nil {
		panic(err)
	}
	vStock := conn.Model(&model.Stock{}).Unscoped().
		Select("stocks.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = stocks.create_by").
		Joins("left join users u2 on u2.id = stocks.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCK, gorm.ViewOption{
		Replace: true,
		Query:   vStock,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCKLOG)
	if err != nil {
		panic(err)
	}
	vStocklog := conn.Model(&model.Stocklog{}).Unscoped().
		Select("stocklogs.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = stocklogs.create_by").
		Joins("left join users u2 on u2.id = stocklogs.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCKLOG, gorm.ViewOption{
		Replace: true,
		Query:   vStocklog,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_VEHICLE)
	if err != nil {
		panic(err)
	}
	vVehicle := conn.Model(&model.Vehicle{}).Unscoped().
		Select("vehicles.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = vehicles.create_by").
		Joins("left join users u2 on u2.id = vehicles.update_by")
	err = conn.Migrator().CreateView(model.VIEW_VEHICLE, gorm.ViewOption{
		Replace: true,
		Query:   vVehicle,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_PRODUCT)
	if err != nil {
		panic(err)
	}
	vProduct := conn.Model(&model.Product{}).Unscoped().
		Select("products.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = products.create_by").
		Joins("left join users u2 on u2.id = products.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PRODUCT, gorm.ViewOption{
		Replace: true,
		Query:   vProduct,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCKMOVEMENT)
	if err != nil {
		panic(err)
	}
	vStockmovement := conn.Model(&model.Stockmovement{}).Unscoped().
		Select("stockmovements.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = stockmovements.create_by").
		Joins("left join users u2 on u2.id = stockmovements.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCKMOVEMENT, gorm.ViewOption{
		Replace: true,
		Query:   vStockmovement,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCKMOVEMENTVEHICLE)
	if err != nil {
		panic(err)
	}
	vStockmovementvehicle := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select([]string{
			"stockmovementvehicles.*, stockmovements.from_warehouse_id, stockmovements.to_warehouse_id, stockmovements.type",
			"CASE WHEN stockmovementvehicles.recived_time IS NOT NULL THEN stockmovementvehicles.sent_net_quantity - stockmovementvehicles.recived_net_quantity ELSE NULL END AS shrinkage",
			"CASE " +
				"WHEN stockmovementvehicles.sent_time IS NULL THEN 'LOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity = 0 THEN 'IN TRANSIT' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity != 0 THEN 'UNLOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NOT NULL THEN 'COMPLETED' " +
				"ELSE '' END AS status",
			"u1.fullname as create_name, u2.fullname as update_name",
		}).
		Joins("left join stockmovements on stockmovements.id = stockmovementvehicles.stockmovement_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCKMOVEMENTVEHICLE, gorm.ViewOption{
		Replace: true,
		Query:   vStockmovementvehicle,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_INBOUND)
	if err != nil {
		panic(err)
	}
	vInbound := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select([]string{
			"stockmovementvehicles.*, stockmovements.to_warehouse_id as warehouse_id, stockmovements.type, stockmovements.remark",
			"CASE WHEN stockmovementvehicles.recived_time IS NOT NULL THEN stockmovementvehicles.sent_net_quantity - stockmovementvehicles.recived_net_quantity ELSE NULL END AS shrinkage",
			"CASE " +
				"WHEN stockmovementvehicles.sent_time IS NULL THEN 'LOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity = 0 THEN 'IN TRANSIT' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity != 0 THEN 'UNLOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NOT NULL THEN 'COMPLETED' " +
				"ELSE '' END AS status",
			"u1.fullname as create_name, u2.fullname as update_name",
		}).
		Joins("join stockmovements stockmovements on stockmovements.id = stockmovementvehicles.stockmovement_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by").
		Where("stockmovements.type = ?", model.StockMovementTypeTransfer).
		Where("stockmovementvehicles.sent_time IS NOT NULL")
	err = conn.Migrator().CreateView(model.VIEW_INBOUND, gorm.ViewOption{
		Replace: true,
		Query:   vInbound,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_OUTBOUND)
	if err != nil {
		panic(err)
	}
	vOutbound := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select([]string{
			"stockmovementvehicles.*, stockmovements.from_warehouse_id as warehouse_id, stockmovements.type, stockmovements.remark",
			"CASE WHEN stockmovementvehicles.recived_time IS NOT NULL THEN stockmovementvehicles.sent_net_quantity - stockmovementvehicles.recived_net_quantity ELSE NULL END AS shrinkage",
			"CASE " +
				"WHEN stockmovementvehicles.sent_time IS NULL THEN 'LOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity = 0 THEN 'IN TRANSIT' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NULL AND stockmovementvehicles.recived_gross_quantity != 0 THEN 'UNLOADING' " +
				"WHEN stockmovementvehicles.sent_time IS NOT NULL AND stockmovementvehicles.recived_time IS NOT NULL THEN 'COMPLETED' " +
				"ELSE '' END AS status",
			"u1.fullname as create_name, u2.fullname as update_name",
		}).
		Joins("join stockmovements stockmovements on stockmovements.id = stockmovementvehicles.stockmovement_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by").
		Where("stockmovements.type = ?", model.StockMovementTypeTransfer)
	err = conn.Migrator().CreateView(model.VIEW_OUTBOUND, gorm.ViewOption{
		Replace: true,
		Query:   vOutbound,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCKIN)
	if err != nil {
		panic(err)
	}
	vStockin := conn.Model(&model.Stockmovement{}).Unscoped().
		Select("stockmovements.*, stockmovements.to_warehouse_id as warehouse_id, stocklogs.gross_quantity, stocklogs.tare_quantity, stocklogs.net_quantity, u1.fullname as create_name, u2.fullname as update_name").
		Joins("join stocklogs stocklogs on stocklogs.stockmovement_id = stockmovements.id").
		Joins("left join users u1 on u1.id = stockmovements.create_by").
		Joins("left join users u2 on u2.id = stockmovements.update_by").
		Where("stockmovements.type = ?", model.StockMovementTypeIn)
	err = conn.Migrator().CreateView(model.VIEW_STOCKIN, gorm.ViewOption{
		Replace: true,
		Query:   vStockin,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_DELIVERYRETAIL)
	if err != nil {
		panic(err)
	}
	vDeliveryretail := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select("stockmovementvehicles.*, stockmovements.from_warehouse_id as warehouse_id, stockmovements.type, stockmovements.remark,  u1.fullname as create_name, u2.fullname as update_name").
		Joins("join stockmovements stockmovements on stockmovements.id = stockmovementvehicles.stockmovement_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by").
		Where("stockmovements.type = ?", model.StockMovementTypeRetail)
	err = conn.Migrator().CreateView(model.VIEW_DELIVERYRETAIL, gorm.ViewOption{
		Replace: true,
		Query:   vDeliveryretail,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_DELIVERYPURCHASEORDER)
	if err != nil {
		panic(err)
	}
	vDeliverypurchaseorder := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select("stockmovementvehicles.*, stockmovements.from_warehouse_id as warehouse_id, stockmovements.type, stockmovements.remark,  u1.fullname as create_name, u2.fullname as update_name").
		Joins("join stockmovements stockmovements on stockmovements.id = stockmovementvehicles.stockmovement_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by").
		Where("stockmovements.type = ?", model.StockMovementTypePurchaseOrder)
	err = conn.Migrator().CreateView(model.VIEW_DELIVERYPURCHASEORDER, gorm.ViewOption{
		Replace: true,
		Query:   vDeliverypurchaseorder,
	})
	if err != nil {
		panic(err)
	}

}

func dbDown() {
	log.Info("Reverting database migrations...")
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Exec("DROP SCHEMA public CASCADE").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("CREATE SCHEMA public").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO postgres").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO public").Error
	if err != nil {
		panic(err)
	}
}

func dbSeed() {
	log.Info("Seeding the database with initial data start")

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	now := time.Now()

	batubara := "1e4f506a-0489-4324-8cd4-005e744d77d0"

	kalimantan := "1492a502-5bf7-4332-a0c7-5e75bba6cc41"
	marunda := "505d34bd-2b2f-4123-a6d2-a56518e2d61f"
	purwakarta := "3aa54082-1f90-49df-b9fa-55b9e85135d2"

	admin := "7aee971f-4e84-4636-aaa2-8dc5fbde2d6b"
	opkalimantan1 := "44427418-9e89-48e4-8ae9-2dabcb56fd84"
	opkalimantan2 := "0db5826f-9640-4324-95ef-5036043ec92a"
	opmarunda1 := "b8171dc6-b9eb-425e-ab75-bf627c7f04d1"
	opmarunda2 := "feb566ea-783e-4fa4-97c6-ea1fa47f8a94"
	oppurwakarta1 := "fad26e34-502b-4b28-9b41-05e81d742b42"
	oppurwakarta2 := "29d834ae-ea1e-4058-8cb1-8ede37424f76"

	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	products := []model.Product{
		{ID: batubara, Name: "Batu Bara"},
	}
	tx.Create(&products)

	users := []model.User{
		{ID: admin, WarehouseID: "", Role: model.UserRoleAdmin, Email: "admin@gmail.com", Username: "admin", PhoneNumber: utils.FormatPhoneTo62("6287770333043"), Fullname: "Admin", Address: "Jl. Gunung Sahari No. 10, Jakarta Pusat", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opkalimantan1, WarehouseID: kalimantan, Role: model.UserRoleOperator, Email: "opkalimantan1@gmail.com", Username: "opkalimantan1", PhoneNumber: "", Fullname: "opkalimantan1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opkalimantan2, WarehouseID: kalimantan, Role: model.UserRoleOperator, Email: "opkalimantan2@gmail.com", Username: "opkalimantan2", PhoneNumber: "", Fullname: "opkalimantan2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opmarunda1, WarehouseID: marunda, Role: model.UserRoleOperator, Email: "opmarunda1@gmail.com", Username: "opmarunda1", PhoneNumber: "", Fullname: "opmarunda1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opmarunda2, WarehouseID: marunda, Role: model.UserRoleOperator, Email: "opmarunda2@gmail.com", Username: "opmarunda2", PhoneNumber: "", Fullname: "opmarunda2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oppurwakarta1, WarehouseID: purwakarta, Role: model.UserRoleOperator, Email: "oppurwakarta1@gmail.com", Username: "oppurwakarta1", PhoneNumber: "", Fullname: "oppurwakarta1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oppurwakarta2, WarehouseID: purwakarta, Role: model.UserRoleOperator, Email: "oppurwakarta2@gmail.com", Username: "oppurwakarta2", PhoneNumber: "", Fullname: "oppurwakarta2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
	}
	tx.Create(&users)

	userprivileges := []model.Userprivilege{
		{UserID: opkalimantan1, StockIn: true, TransferOut: true, TransferIn: true, PurchaseOrder: true, Retail: true},
		{UserID: opkalimantan2, StockIn: false, TransferOut: true, TransferIn: false, PurchaseOrder: false, Retail: false},
		{UserID: opmarunda1, StockIn: true, TransferOut: true, TransferIn: true, PurchaseOrder: true, Retail: true},
		{UserID: opmarunda2, StockIn: false, TransferOut: true, TransferIn: false, PurchaseOrder: true, Retail: true},
		{UserID: oppurwakarta1, StockIn: true, TransferOut: true, TransferIn: true, PurchaseOrder: true, Retail: true},
		{UserID: oppurwakarta2, StockIn: false, TransferOut: false, TransferIn: true, PurchaseOrder: true, Retail: true},
	}
	tx.Create(&userprivileges)

	warehouses := []model.Warehouse{
		{ID: kalimantan, Name: "Kalimantan", Location: "Kalimantan", IsStockin: true, IsInbound: false, IsOutbound: true, IsRetail: false, IsPurchaseorder: false},
		{ID: marunda, Name: "Marunda", Location: "Marunda", IsStockin: false, IsInbound: true, IsOutbound: true, IsRetail: false, IsPurchaseorder: false},
		{ID: purwakarta, Name: "Purwakarta", Location: "Purwakarta", IsStockin: false, IsInbound: true, IsOutbound: false, IsRetail: true, IsPurchaseorder: true},
	}
	tx.Create(&warehouses)

	stocks := []model.Stock{
		{WarehouseID: kalimantan, ProductID: batubara, Quantity: 0},
		{WarehouseID: marunda, ProductID: batubara, Quantity: 0},
		{WarehouseID: purwakarta, ProductID: batubara, Quantity: 0},
	}
	tx.Create(&stocks)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}

	log.Info("Seeding the database with initial data end")
}

func dbReset() {
	dbDown()
	dbUp()
	dbSeed()
}
