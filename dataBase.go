package main

import (
	"Taskl0/model"
	"database/sql"
	_ "encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/patrickmn/go-cache"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "skar4500"
	dbname   = "postgres"
)

type DB struct {
	conn *sql.DB
}

func (db *DB) addOrderDataDB(order model.OrderModel) error {
	fmt.Println("add order data in data base")
	psqlconn := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=disable", host, port, dbname, user, password)

	db.conn, _ = sql.Open("postgres", psqlconn)
	defer db.conn.Close()

	//var query string = "CALL public.\"addOrder\"('" + idData + "','" + string(data) + "');"

	query := `INSERT INTO wb_order(
                   order_uid, track_number, entry, locale,
                   internal_signature, customer_id, delivery_service, shardkey,
                   sm_id, date_created, oof_shard
                   ) 
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	err := db.conn.QueryRow(query, order.OrderUID, order.TrackNumber, order.Entry, order.Local,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard).Err()
	if err != nil {
		return err
	}

	query = `INSERT INTO delivery(name, phone, zip, city, address, region, email, order_uid)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	err = db.conn.QueryRow(query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
		order.OrderUID).Err()
	if err != nil {
		return err
	}

	query = `INSERT INTO payment(transaction, request_id, currency, provider, amount,
		payment_dt, bank, delivery_cost, goods_total, custom_fee, order_uid)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	err = db.conn.QueryRow(query, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID).Err()
	if err != nil {
		return err
	}

	query = `INSERT INTO items(chrt_id, track_number, price, rid, name,
                  sale, size, total_price, nm_id, brand, status, order_uid) 
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	err = db.conn.QueryRow(query, order.Items[0].ChrtId, order.Items[0].TrackNumber, order.Items[0].Price,
		order.Items[0].Rid, order.Items[0].Name, order.Items[0].Sale, order.Items[0].Size,
		order.Items[0].TotalPrice, order.Items[0].NmID, order.Items[0].Brand,
		order.Items[0].Status, order.OrderUID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getOrderDataDB(idData string) model.OrderModel {
	psqlconn := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=disable", host, port, dbname, user, password)

	var err error
	db.conn, err = sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.conn.Close()

	//var query string = "SELECT public.\"getDateId\"('" + idData + "');"

	query := `SELECT wb_order.order_uid, wb_order.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	  delivery.name, phone, zip, city, address, region, email,
	  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee,
	  chrt_id, items.track_number, price, rid, items.name, sale, size, total_price, nm_id, brand, status
		  FROM wb_order, delivery, payment, items
		  WHERE wb_order.order_uid = $1 AND 
		        delivery.order_uid = wb_order.order_uid AND
		        payment.order_uid = wb_order.order_uid AND
		        items.order_uid = wb_order.order_uid`

	var result model.OrderModel

	db.conn.QueryRow(query, idData).Scan(&result.OrderUID, &result.TrackNumber, &result.Entry, &result.Local, &result.InternalSignature,
		&result.CustomerID, &result.DeliveryService, &result.Shardkey, &result.SmID, &result.DateCreated, &result.OofShard,
		&result.Delivery.Name, &result.Delivery.Phone, &result.Delivery.Zip, &result.Delivery.City, &result.Delivery.Address,
		&result.Delivery.Region, &result.Delivery.Email, &result.Payment.Transaction, &result.Payment.RequestID, &result.Payment.Currency,
		&result.Payment.Provider, &result.Payment.Amount, &result.Payment.PaymentDt, &result.Payment.Bank, &result.Payment.DeliveryCost,
		&result.Payment.GoodsTotal, &result.Payment.CustomFee, &result.Items[0].ChrtId, &result.Items[0].TrackNumber, &result.Items[0].Price,
		&result.Items[0].Rid, &result.Items[0].Name, &result.Items[0].Sale, &result.Items[0].Size, &result.Items[0].TotalPrice,
		&result.Items[0].NmID, &result.Items[0].Brand, &result.Items[0].Status)

	return result
}

func (db *DB) getAllOrdersDataDB() ([]model.OrderModel, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=disable", host, port, dbname, user, password)

	var err error
	db.conn, err = sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.conn.Close()

	query := `SELECT wb_order.order_uid, wb_order.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	  delivery.name, phone, zip, city, address, region, email,
	  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee,
	  chrt_id, items.track_number, price, rid, items.name, sale, size, total_price, nm_id, brand, status
		  FROM wb_order, delivery, payment, items
		  WHERE items.order_uid = wb_order.order_uid AND payment.order_uid = wb_order.order_uid AND delivery.order_uid = wb_order.order_uid`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	wb_orders := make([]model.OrderModel, 0)
	for rows.Next() {
		var order model.OrderModel
		order.Items = make([]model.ItemModel, 1)
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Local, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee, &order.Items[0].ChrtId, &order.Items[0].TrackNumber, &order.Items[0].Price,
			&order.Items[0].Rid, &order.Items[0].Name, &order.Items[0].Sale, &order.Items[0].Size, &order.Items[0].TotalPrice,
			&order.Items[0].NmID, &order.Items[0].Brand, &order.Items[0].Status)
		if err != nil {
			return nil, err
		}
		wb_orders = append(wb_orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wb_orders, nil
}

func addErrData(errDataString string) {
	fmt.Println("add order err in data base")
	psqlconn := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=disable", host, port, dbname, user, password)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var query string = "CALL public.\"addErrData\"('" + errDataString + "');"

	db.Exec(query)
}

//func OrdersDBCache(cacheProgram *cache.Cache) *cache.Cache {
//
//	fmt.Println("restart cache = db")
//	psqlconn := fmt.Sprintf("host=%s port=%d dbname=%s user='%s' password=%s sslmode=disable", host, port, dbname, user, password)
//	db, err := sql.Open("postgres", psqlconn)
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer db.Close()
//
//	//var query string =
//
//	//var query string = "SELECT public.\"getOrders\"();"
//	var query string = "" //"SELECT public.\"getOrders\"();"
//	row, err := db.Query(query)
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer row.Close()
//
//	for row.Next() {
//		var jsonOrderDB []byte
//		var order map[string]interface{}
//
//		row.Scan(&jsonOrderDB)
//
//		json.Unmarshal(jsonOrderDB, &order)
//
//		cacheProgram.Add(fmt.Sprintf("%+v", order["order_uid"]), string(jsonOrderDB), cache.DefaultExpiration)
//
//	}
//
//	cacheProgram.SaveFile("cache")
//
//	return cacheProgram
//
//}
