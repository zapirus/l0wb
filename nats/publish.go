package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type order struct {
	Space string
	Point json.RawMessage
}

func PublishMsg(msg []byte) {
	sc, err := stan.Connect("test-cluster", "pub", stan.NatsURL("localhost:4222"))
	if err != nil {
		fmt.Printf("An error occures in publish func: %v\n", err)
		return
	}
	defer sc.Close()
	sc.Publish("test-channel", msg)
	fmt.Printf("Published msg: %s\n", msg)
}

func main() {
	// start publishing
	go func() {
		for {
			PublishMsg([]byte(`{
			"order_uid": "` + fmt.Sprintf("%d", rand.Int()%1000) + `",
			"track_number": "WBILMTESTTRACK",
			"entry": "WBIL",
			"delivery": {
			  "name": "Test1 Testov",
			  "phone": "+9720000000",
			  "zip": "2639809",
			  "city": "Kiryat Mozkin",
			  "address": "Ploshad Mira 15",
			  "region": "Kraiot",
			  "email": "test@gmail.com"
			},
			"payment": {
			  "transaction": "` + fmt.Sprintf("%d", rand.Int()%1000) + `",
			  "request_id": "",
			  "currency": "USD",
			  "provider": "wbpay",
			  "amount": 1817,
			  "payment_dt": 1637907727,
			  "bank": "alpha",
			  "delivery_cost": 1500,
			  "goods_total": 317,
			  "custom_fee": 0
			},
			"items": [
			  {
				"chrt_id": ` + fmt.Sprintf("%d", rand.Int()%100000) + `,
				"track_number": "WBILMTESTTRACK",
				"price": ` + fmt.Sprintf("%d", rand.Int()%1000) + `,
				"rid": "a",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			  }
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "test",
			"delivery_service": "meest",
			"shardkey": "9",
			"sm_id": 99,
			"date_created": "2021-11-26T06:22:19Z",
			"oof_shard": "1"
		  }`))
			time.Sleep(time.Second * 25)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Exit by signal")
}
