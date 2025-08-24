package cmd

import (
	"time"

	"github.com/jihanlugas/warehouse/cryption"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
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

	err = conn.Migrator().AutoMigrate(&model.Photo{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Photoinc{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Userprovider{})
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
	err = conn.Migrator().AutoMigrate(&model.Retailproduct{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Purchaseorder{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Purchaseorderproduct{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Transaction{})
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
	err = conn.Migrator().AutoMigrate(&model.Location{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Warehouse{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Warehousedestination{})
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
	err = conn.Migrator().AutoMigrate(&model.Stockmovementvehicle{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Stockmovementvehiclephoto{})
	if err != nil {
		panic(err)
	}
}

func dbUpView() {
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Migrator().DropView(model.VIEW_PHOTO)
	if err != nil {
		panic(err)
	}
	vPhoto := conn.Model(&model.Photo{}).Unscoped().
		Select("photos.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = photos.create_by").
		Joins("left join users u2 on u2.id = photos.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PHOTO, gorm.ViewOption{
		Replace: true,
		Query:   vPhoto,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_PHOTOINC)
	if err != nil {
		panic(err)
	}
	vPhotoinc := conn.Model(&model.Photoinc{}).Unscoped().
		Select("photoincs.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = photoincs.create_by").
		Joins("left join users u2 on u2.id = photoincs.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PHOTOINC, gorm.ViewOption{
		Replace: true,
		Query:   vPhotoinc,
	})
	if err != nil {
		panic(err)
	}

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

	err = conn.Migrator().DropView(model.VIEW_USERPROVIDER)
	if err != nil {
		panic(err)
	}
	vUserprovider := conn.Model(&model.Userprovider{}).Unscoped().
		Select("userproviders.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = userproviders.create_by").
		Joins("left join users u2 on u2.id = userproviders.update_by")
	err = conn.Migrator().CreateView(model.VIEW_USERPROVIDER, gorm.ViewOption{
		Replace: true,
		Query:   vUserprovider,
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
		Select("retails.*" +
			", coalesce(stockmovementvehicles.total_price, 0) as total_price" +
			", coalesce(transactions.total_payment, 0) as total_payment" +
			", coalesce((total_price - coalesce(total_payment, 0)), 0) as outstanding " +
			", u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join ( " +
			"select stockmovementvehicles.related_id, coalesce(sum(retailproducts.unit_price * stockmovementvehicles.sent_net_quantity), 0) as total_price " +
			"from stockmovementvehicles stockmovementvehicles " +
			"join retails retails on retails.id = stockmovementvehicles.related_id " +
			"join retailproducts retailproducts on retailproducts.retail_id = retails.id and retailproducts.product_id = stockmovementvehicles.product_id " +
			"where stockmovementvehicles.delete_dt is null " +
			"and stockmovementvehicles.sent_time is not null " +
			"group by stockmovementvehicles.related_id " +
			" ) as stockmovementvehicles on stockmovementvehicles.related_id = retails.id").
		Joins("left join ( " +
			"select transactions.related_id, coalesce(sum(transactions.amount), 0) as total_payment " +
			"from transactions transactions join retails on retails.id = transactions.related_id " +
			"where retails.delete_dt is null " +
			"group by transactions.related_id " +
			") as transactions on transactions.related_id = retails.id").
		Joins("left join users u1 on u1.id = retails.create_by").
		Joins("left join users u2 on u2.id = retails.update_by")
	err = conn.Migrator().CreateView(model.VIEW_RETAIL, gorm.ViewOption{
		Replace: true,
		Query:   vRetail,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_RETAILPRODUCT)
	if err != nil {
		panic(err)
	}
	vRetailproduct := conn.Model(&model.Retailproduct{}).Unscoped().
		Select("retailproducts.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = retailproducts.create_by").
		Joins("left join users u2 on u2.id = retailproducts.update_by")
	err = conn.Migrator().CreateView(model.VIEW_RETAILPRODUCT, gorm.ViewOption{
		Replace: true,
		Query:   vRetailproduct,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_PURCHASEORDER)
	if err != nil {
		panic(err)
	}
	vPurchaseorder := conn.Model(&model.Purchaseorder{}).Unscoped().
		Select("purchaseorders.*" +
			", coalesce(stockmovementvehicles.total_price, 0) as total_price" +
			", coalesce(transactions.total_payment, 0) as total_payment" +
			", coalesce((total_price - coalesce(total_payment, 0)), 0) as outstanding " +
			", u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join ( " +
			"select stockmovementvehicles.related_id, coalesce(sum(purchaseorderproducts.unit_price * stockmovementvehicles.sent_net_quantity), 0) as total_price " +
			"from stockmovementvehicles stockmovementvehicles " +
			"join purchaseorders purchaseorders on purchaseorders.id = stockmovementvehicles.related_id " +
			"join purchaseorderproducts purchaseorderproducts on purchaseorderproducts.purchaseorder_id = purchaseorders.id and purchaseorderproducts.product_id = stockmovementvehicles.product_id " +
			"where stockmovementvehicles.delete_dt is null " +
			"and stockmovementvehicles.sent_time is not null " +
			"group by stockmovementvehicles.related_id " +
			" ) as stockmovementvehicles on stockmovementvehicles.related_id = purchaseorders.id").
		Joins("left join ( " +
			"select transactions.related_id, coalesce(sum(transactions.amount), 0) as total_payment " +
			"from transactions transactions join purchaseorders on purchaseorders.id = transactions.related_id " +
			"where purchaseorders.delete_dt is null " +
			"group by transactions.related_id " +
			") as transactions on transactions.related_id = purchaseorders.id").
		Joins("left join users u1 on u1.id = purchaseorders.create_by").
		Joins("left join users u2 on u2.id = purchaseorders.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PURCHASEORDER, gorm.ViewOption{
		Replace: true,
		Query:   vPurchaseorder,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_PURCHASEORDERPRODUCT)
	if err != nil {
		panic(err)
	}
	vPurchaseorderproduct := conn.Model(&model.Purchaseorderproduct{}).Unscoped().
		Select("purchaseorderproducts.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = purchaseorderproducts.create_by").
		Joins("left join users u2 on u2.id = purchaseorderproducts.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PURCHASEORDERPRODUCT, gorm.ViewOption{
		Replace: true,
		Query:   vPurchaseorderproduct,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_TRANSACTION)
	if err != nil {
		panic(err)
	}
	vTransaction := conn.Model(&model.Transaction{}).Unscoped().
		Select([]string{
			"transactions.*, u1.fullname as create_name, u2.fullname as update_name",
			"CASE " +
				"WHEN transactions.transaction_related = 'PURCHASE_ORDER' THEN purchaseorders.customer_id " +
				"WHEN transactions.transaction_related = 'RETAIL' THEN retails.customer_id " +
				"ELSE '' END AS customer_id",
		}).
		Joins("left join purchaseorders purchaseorders on purchaseorders.id = transactions.related_id").
		Joins("left join retails retails on retails.id = transactions.related_id").
		Joins("left join users u1 on u1.id = transactions.create_by").
		Joins("left join users u2 on u2.id = transactions.update_by")
	err = conn.Migrator().CreateView(model.VIEW_TRANSACTION, gorm.ViewOption{
		Replace: true,
		Query:   vTransaction,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_LOCATION)
	if err != nil {
		panic(err)
	}
	vLocation := conn.Model(&model.Location{}).Unscoped().
		Select("locations.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = locations.create_by").
		Joins("left join users u2 on u2.id = locations.update_by")
	err = conn.Migrator().CreateView(model.VIEW_LOCATION, gorm.ViewOption{
		Replace: true,
		Query:   vLocation,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_WAREHOUSE)
	if err != nil {
		panic(err)
	}
	vWarehouse := conn.Model(&model.Warehouse{}).Unscoped().
		Select("warehouses.*, coalesce(transferouts.total_running_transferout, 0) as total_running_transferout, coalesce(transferins.total_running_transferin, 0) as total_running_transferin, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join ( " +
			"select warehouses.id, count(warehouses.id) as total_running_transferout " +
			"from stockmovementvehicles stockmovementvehicles " +
			"join warehouses warehouses on warehouses.id = stockmovementvehicles.from_warehouse_id " +
			"where stockmovementvehicles.stockmovementvehicle_status != 'LOADING' " +
			"and stockmovementvehicles.stockmovementvehicle_status != 'COMPLETED'  " +
			"and stockmovementvehicles.delete_dt is null " +
			"and stockmovementvehicles.stockmovementvehicle_type = 'TRANSFER' " +
			"and stockmovementvehicles.delete_dt is null " +
			"group by warehouses.id " +
			") as transferouts on transferouts.id = warehouses.id").
		Joins("left join ( " +
			"select warehouses.id, count(warehouses.id) as total_running_transferin " +
			"from stockmovementvehicles stockmovementvehicles " +
			"join warehouses warehouses on warehouses.id = stockmovementvehicles.to_warehouse_id " +
			"where stockmovementvehicles.stockmovementvehicle_status != 'LOADING' " +
			"and stockmovementvehicles.stockmovementvehicle_status != 'COMPLETED'  " +
			"and stockmovementvehicles.delete_dt is null " +
			"and stockmovementvehicles.stockmovementvehicle_type = 'TRANSFER' " +
			"and stockmovementvehicles.delete_dt is null " +
			"group by warehouses.id " +
			") as transferins on transferins.id = warehouses.id").
		Joins("left join users u1 on u1.id = warehouses.create_by").
		Joins("left join users u2 on u2.id = warehouses.update_by")
	err = conn.Migrator().CreateView(model.VIEW_WAREHOUSE, gorm.ViewOption{
		Replace: true,
		Query:   vWarehouse,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_WAREHOUSEDESTINATION)
	if err != nil {
		panic(err)
	}
	vWarehousedestination := conn.Model(&model.Warehousedestination{}).Unscoped().
		Select("warehousedestinations.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = warehousedestinations.create_by").
		Joins("left join users u2 on u2.id = warehousedestinations.update_by")
	err = conn.Migrator().CreateView(model.VIEW_WAREHOUSEDESTINATION, gorm.ViewOption{
		Replace: true,
		Query:   vWarehousedestination,
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

	err = conn.Migrator().DropView(model.VIEW_STOCKMOVEMENTVEHICLE)
	if err != nil {
		panic(err)
	}
	vStockmovementvehicle := conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Select([]string{
			"stockmovementvehicles.*",
			"CASE " +
				"WHEN stockmovementvehicles.stockmovementvehicle_type = 'PURCHASE_ORDER' THEN coalesce(purchaseorderproducts.unit_price, 0) " +
				"WHEN stockmovementvehicles.stockmovementvehicle_type = 'RETAIL' THEN coalesce(retailproducts.unit_price, 0) " +
				"ELSE 0 END AS unit_price ",
			"u1.fullname as create_name, u2.fullname as update_name",
		}).
		Joins("left join products products on products.id = stockmovementvehicles.product_id").
		Joins("left join purchaseorders purchaseorders on purchaseorders.id = stockmovementvehicles.related_id").
		Joins("left join purchaseorderproducts purchaseorderproducts on purchaseorderproducts.purchaseorder_id = purchaseorders.id and stockmovementvehicles.product_id = purchaseorderproducts.product_id").
		Joins("left join retails retails on retails.id = stockmovementvehicles.related_id").
		Joins("left join retailproducts retailproducts on retailproducts.retail_id = retails.id and stockmovementvehicles.product_id = retailproducts.product_id").
		Joins("left join users u1 on u1.id = stockmovementvehicles.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehicles.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCKMOVEMENTVEHICLE, gorm.ViewOption{
		Replace: true,
		Query:   vStockmovementvehicle,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_STOCKMOVEMENTVEHICLEPHOTO)
	if err != nil {
		panic(err)
	}
	vStockmovementvehiclephoto := conn.Model(&model.Stockmovementvehiclephoto{}).Unscoped().
		Select("stockmovementvehiclephotos.*, photos.photo_path as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join photos photos on photos.id = stockmovementvehiclephotos.photo_id").
		Joins("left join users u1 on u1.id = stockmovementvehiclephotos.create_by").
		Joins("left join users u2 on u2.id = stockmovementvehiclephotos.update_by")
	err = conn.Migrator().CreateView(model.VIEW_STOCKMOVEMENTVEHICLEPHOTO, gorm.ViewOption{
		Replace: true,
		Query:   vStockmovementvehiclephoto,
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

	// Locations
	locationkaltim := "1492a502-5bf7-4332-a0c7-5e75bba6cc41"
	locationkalsel := "66b4de0d-2dce-497f-aac6-912960421a9b"
	locationmarunda := "505d34bd-2b2f-4123-a6d2-a56518e2d61f"
	locationpurwakarta := "3aa54082-1f90-49df-b9fa-55b9e85135d2"

	// Warehouse Kaltim
	rkaltim1 := "66b4de0d-2dce-497f-aac6-912960421a9b"
	rkaltim2 := "cfeb843e-138f-490f-a048-9614c69ef23d"
	rkaltim3 := "720405fa-be79-4649-a72e-439c3cc7c4c6"
	rkaltim4 := "27fbcb79-7139-4261-a03f-f0131b4081c3"
	jkaltim1 := "89ea0924-8200-4337-8716-5c3d7e00df24"
	jkaltim2 := "d5bebb90-f7c3-47ed-8f50-2d80360a65dd"

	// Warehouse Kalsel
	rkalsel1 := "eb259cf4-9643-457d-a9b5-09ff449a0432"
	rkalsel2 := "4c10fa17-f81b-421d-9f09-d6cb0cefec71"
	jkalsel1 := "9da13b12-e151-4792-ab82-3950dc2a6ce9"

	// Warehouse Marunda
	jmarunda1 := "108f3e98-ff43-4bda-9005-213dc90da642"
	smarunda1 := "cd17dda0-4736-4904-9609-24b0ba31c068"

	// Warehouse Purwakarta
	spurwakarta1 := "4a770500-4797-488f-af06-ea77c42110d5"

	admin := "7aee971f-4e84-4636-aaa2-8dc5fbde2d6b"
	oprkaltim1 := "5c777011-3f7f-478b-93ed-2d9ecd804c46"
	oprkaltim2 := "44427418-9e89-48e4-8ae9-2dabcb56fd84"
	oprkaltim3 := "0db5826f-9640-4324-95ef-5036043ec92a"
	oprkaltim4 := "27833752-8855-4d73-a4fe-f8d58876ecda"
	opjkaltim1 := "b8171dc6-b9eb-425e-ab75-bf627c7f04d1"
	opjkaltim2 := "feb566ea-783e-4fa4-97c6-ea1fa47f8a94"
	oprkalsel1 := "a4f5c345-8862-4fc8-842f-87194c2cbfcc"
	oprkalsel2 := "fad26e34-502b-4b28-9b41-05e81d742b42"
	opjkalsel1 := "29d834ae-ea1e-4058-8cb1-8ede37424f76"
	opjmarunda1 := "73e13772-68bc-44a5-a381-9cc9975c096c"
	opsmarunda1 := "80e8e449-a724-432f-84c6-d4d345dac46c"
	opspurwakarta1 := "3c8fe8ac-684c-46bf-bf9a-54c9d078383d"

	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	products := []model.Product{
		{ID: batubara, Name: "Batu Bara"},
	}
	tx.Create(&products)

	users := []model.User{
		{ID: admin, LocationID: "", WarehouseID: "", UserRole: model.UserRoleAdmin, Email: "admin@gmail.com", Username: "admin", PhoneNumber: utils.FormatPhoneTo62("6287770333043"), Fullname: "Admin", Address: "Jl. Gunung Sahari No. 10, Jakarta Pusat", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkaltim1, LocationID: locationkaltim, WarehouseID: rkaltim1, UserRole: model.UserRoleOperator, Email: "oprkaltim1@gmail.com", Username: "oprkaltim1", PhoneNumber: "", Fullname: "OP rkaltim1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkaltim2, LocationID: locationkaltim, WarehouseID: rkaltim2, UserRole: model.UserRoleOperator, Email: "oprkaltim2@gmail.com", Username: "oprkaltim2", PhoneNumber: "", Fullname: "OP rkaltim2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkaltim3, LocationID: locationkaltim, WarehouseID: rkaltim3, UserRole: model.UserRoleOperator, Email: "oprkaltim3@gmail.com", Username: "oprkaltim3", PhoneNumber: "", Fullname: "OP rkaltim3", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkaltim4, LocationID: locationkaltim, WarehouseID: rkaltim4, UserRole: model.UserRoleOperator, Email: "oprkaltim4@gmail.com", Username: "oprkaltim4", PhoneNumber: "", Fullname: "OP rkaltim4", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opjkaltim1, LocationID: locationkaltim, WarehouseID: jkaltim1, UserRole: model.UserRoleOperator, Email: "opjkaltim1@gmail.com", Username: "opjkaltim1", PhoneNumber: "", Fullname: "OP jkaltim1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opjkaltim2, LocationID: locationkaltim, WarehouseID: jkaltim2, UserRole: model.UserRoleOperator, Email: "opjkaltim2@gmail.com", Username: "opjkaltim2", PhoneNumber: "", Fullname: "OP jkaltim2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkalsel1, LocationID: locationkalsel, WarehouseID: rkalsel1, UserRole: model.UserRoleOperator, Email: "oprkalsel1@gmail.com", Username: "oprkalsel1", PhoneNumber: "", Fullname: "OP rkalsel1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: oprkalsel2, LocationID: locationkalsel, WarehouseID: rkalsel2, UserRole: model.UserRoleOperator, Email: "oprkalsel2@gmail.com", Username: "oprkalsel2", PhoneNumber: "", Fullname: "OP rkalsel2", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opjkalsel1, LocationID: locationkalsel, WarehouseID: jkalsel1, UserRole: model.UserRoleOperator, Email: "opjkalsel1@gmail.com", Username: "opjkalsel1", PhoneNumber: "", Fullname: "OP jkalsel1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opjmarunda1, LocationID: locationmarunda, WarehouseID: jmarunda1, UserRole: model.UserRoleOperator, Email: "opjmarunda1@gmail.com", Username: "opjmarunda1", PhoneNumber: "", Fullname: "OP jmarunda1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opsmarunda1, LocationID: locationmarunda, WarehouseID: smarunda1, UserRole: model.UserRoleOperator, Email: "opsmarunda1@gmail.com", Username: "opsmarunda1", PhoneNumber: "", Fullname: "OP smarunda1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
		{ID: opspurwakarta1, LocationID: locationpurwakarta, WarehouseID: spurwakarta1, UserRole: model.UserRoleOperator, Email: "opspurwakarta1@gmail.com", Username: "opspurwakarta1", PhoneNumber: "", Fullname: "OP spurwakarta1", Address: "", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now},
	}
	tx.Create(&users)

	userprivileges := []model.Userprivilege{
		{UserID: oprkaltim1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: oprkaltim2, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: oprkaltim3, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: oprkaltim4, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opjkaltim1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opjkaltim2, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: oprkalsel1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: oprkalsel2, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opjkalsel1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opjmarunda1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opsmarunda1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
		{UserID: opspurwakarta1, StockIn: true, TransferOut: true, TransferIn: true, Purchaseorder: true, Retail: true},
	}
	tx.Create(&userprivileges)

	locations := []model.Location{
		{ID: locationkalsel, Name: "Kalimantan Selatan", Notes: "Notes Kalimantan Selatan"},
		{ID: locationkaltim, Name: "Kalimantan Timur", Notes: "Notes Kalimantan Timur"},
		{ID: locationmarunda, Name: "Marunda", Notes: "Notes Marunda"},
		{ID: locationpurwakarta, Name: "Purwakarta", Notes: "Notes Purwakarta"},
	}
	tx.Create(&locations)

	warehouses := []model.Warehouse{
		{ID: rkaltim1, LocationID: locationkaltim, Name: "Stockroom Kaltim 1", Address: "Jl. Stockroom Kaltim 1", Notes: "Notes Stockroom Kaltim 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: rkaltim2, LocationID: locationkaltim, Name: "Stockroom Kaltim 2", Address: "Jl. Stockroom Kaltim 2", Notes: "Notes Stockroom Kaltim 2", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: rkaltim3, LocationID: locationkaltim, Name: "Stockroom Kaltim 3", Address: "Jl. Stockroom Kaltim 3", Notes: "Notes Stockroom Kaltim 3", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: rkaltim4, LocationID: locationkaltim, Name: "Stockroom Kaltim 4", Address: "Jl. Stockroom Kaltim 4", Notes: "Notes Stockroom Kaltim 4", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: jkaltim1, LocationID: locationkaltim, Name: "Jetty Kaltim 1", Address: "Jl. Jetty Kaltim 1", Notes: "Notes Jetty Kaltim 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: jkaltim2, LocationID: locationkaltim, Name: "Jetty Kaltim 2", Address: "Jl. Jetty Kaltim 2", Notes: "Notes Jetty Kaltim 2", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: rkalsel1, LocationID: locationkalsel, Name: "Stockroom Kalsel 1", Address: "Jl. Stockroom Kalsel 1", Notes: "Notes Stockroom Kalsel 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: rkalsel2, LocationID: locationkalsel, Name: "Stockroom Kalsel 2", Address: "Jl. Stockroom Kalsel 2", Notes: "Notes Stockroom Kalsel 2", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: true, IsTransferIn: false, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: jkalsel1, LocationID: locationkalsel, Name: "Jetty Kalsel 1", Address: "Jl. Jetty Kalsel 1", Notes: "Notes Jetty Kaltim 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: true, IsRetail: false, IsPurchaseorder: false},
		{ID: jmarunda1, LocationID: locationmarunda, Name: "Jetty Marunda 1", Address: "Jl. Jetty Marunda 1", Notes: "Notes Jetty Marunda 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: true, IsRetail: true, IsPurchaseorder: true},
		{ID: smarunda1, LocationID: locationmarunda, Name: "Stockpile Marunda 1", Address: "Jl. Stockpile Marunda 1", Notes: "Notes Stockpile Marunda 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: true, IsRetail: true, IsPurchaseorder: true},
		{ID: spurwakarta1, LocationID: locationpurwakarta, Name: "Stockpile Purwakarta 1", Address: "Jl. Stockpile Purwakarta 1", Notes: "Notes Stockpile Purwakarta 1", PhoneNumber: utils.FormatPhoneTo62("081234561234"), IsStockin: false, IsTransferIn: true, IsTransferOut: false, IsRetail: true, IsPurchaseorder: true},
	}
	tx.Create(&warehouses)

	warehousedestinations := []model.Warehousedestination{
		{FromLocationID: locationkaltim, ToLocationID: locationkaltim, FromWarehouseID: rkaltim1, ToWarehouseID: jkaltim1},
		{FromLocationID: locationkaltim, ToLocationID: locationkaltim, FromWarehouseID: rkaltim2, ToWarehouseID: jkaltim1},
		{FromLocationID: locationkaltim, ToLocationID: locationkaltim, FromWarehouseID: rkaltim3, ToWarehouseID: jkaltim2},
		{FromLocationID: locationkaltim, ToLocationID: locationkaltim, FromWarehouseID: rkaltim4, ToWarehouseID: jkaltim2},
		{FromLocationID: locationkaltim, ToLocationID: locationmarunda, FromWarehouseID: jkaltim1, ToWarehouseID: jmarunda1},
		{FromLocationID: locationkaltim, ToLocationID: locationmarunda, FromWarehouseID: jkaltim2, ToWarehouseID: jmarunda1},
		{FromLocationID: locationkalsel, ToLocationID: locationkalsel, FromWarehouseID: rkalsel1, ToWarehouseID: jkalsel1},
		{FromLocationID: locationkalsel, ToLocationID: locationkalsel, FromWarehouseID: rkalsel2, ToWarehouseID: jkalsel1},
		{FromLocationID: locationkalsel, ToLocationID: locationmarunda, FromWarehouseID: jkalsel1, ToWarehouseID: jmarunda1},
		{FromLocationID: locationmarunda, ToLocationID: locationmarunda, FromWarehouseID: jmarunda1, ToWarehouseID: smarunda1},
		{FromLocationID: locationmarunda, ToLocationID: locationpurwakarta, FromWarehouseID: jmarunda1, ToWarehouseID: spurwakarta1},
		{FromLocationID: locationmarunda, ToLocationID: locationpurwakarta, FromWarehouseID: smarunda1, ToWarehouseID: spurwakarta1},
	}
	tx.Create(&warehousedestinations)

	stocks := []model.Stock{
		{LocationID: locationkaltim, WarehouseID: rkaltim1, ProductID: batubara, Quantity: 0},
		{LocationID: locationkaltim, WarehouseID: rkaltim2, ProductID: batubara, Quantity: 0},
		{LocationID: locationkaltim, WarehouseID: rkaltim3, ProductID: batubara, Quantity: 0},
		{LocationID: locationkaltim, WarehouseID: rkaltim4, ProductID: batubara, Quantity: 0},
		{LocationID: locationkaltim, WarehouseID: jkaltim1, ProductID: batubara, Quantity: 0},
		{LocationID: locationkaltim, WarehouseID: jkaltim2, ProductID: batubara, Quantity: 0},
		{LocationID: locationkalsel, WarehouseID: rkalsel1, ProductID: batubara, Quantity: 0},
		{LocationID: locationkalsel, WarehouseID: rkalsel2, ProductID: batubara, Quantity: 0},
		{LocationID: locationkalsel, WarehouseID: jkalsel1, ProductID: batubara, Quantity: 0},
		{LocationID: locationmarunda, WarehouseID: jmarunda1, ProductID: batubara, Quantity: 0},
		{LocationID: locationmarunda, WarehouseID: smarunda1, ProductID: batubara, Quantity: 0},
		{LocationID: locationpurwakarta, WarehouseID: spurwakarta1, ProductID: batubara, Quantity: 0},
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
