package main

import (
	"Taskl0/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
	"os/signal"
	"time"
)

func main() {

	//cacheProgram = restartCache(cacheProgram)
	service := MakeCacheService()
	service.Init()

	go subChan(&service)

	go serverHtmlStart(&service)
	time.Sleep(5 * time.Second)

	fmt.Scanln()
}

func subChan(service *CacheService) {
	nc, err := stan.Connect("test-cluster", "sub")
	if err != nil {
		println("not connect channel")
		time.Sleep(5 * time.Second)
		subChan(service)
	}
	fmt.Println("Connected  to channel")

	handl := func(msg *stan.Msg) {
		var order model.OrderModel
		err = json.Unmarshal(msg.Data, &order)

		if err == nil {
			fmt.Println("Adding order to cache and db")
			ok := service.AddOrder(order)
			if !ok {
				fmt.Printf("error create an order %v\n", ok)
			}
		} else {
			fmt.Printf("invalid msg from channel %s\n", msg.Data)
		}
	}

	sub, err := nc.QueueSubscribe("test-channel", "service", handl)
	if err != nil {
		fmt.Printf("Cannot subscribe to channel %s\n", "test-channel")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Exit by signal")

	sub.Unsubscribe()
	nc.Close()
}
