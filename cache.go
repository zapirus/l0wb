package main

import (
	"Taskl0/model"
	"fmt"
	"sync"
)

type Cache struct {
	data map[string]model.OrderModel
	m    sync.RWMutex
}

func New() *Cache {
	return &Cache{
		data: make(map[string]model.OrderModel),
		m:    sync.RWMutex{},
	}
}

func MakeCache() Cache {
	return Cache{
		data: make(map[string]model.OrderModel),
		m:    sync.RWMutex{},
	}
}

func (c *Cache) GetOrderById(id string) (model.OrderModel, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	if o, ok := c.data[id]; ok {
		return o, nil
	}
	return model.OrderModel{}, fmt.Errorf("no order with id: %s", id)
}

func (c *Cache) AddOrder(order model.OrderModel) {
	c.m.Lock()
	defer c.m.Unlock()
	c.data[order.OrderUID] = order
}

//func addDataCache(cacheProgram *cache.Cache, orderId string, orderDataJson []byte, mutex *sync.Mutex) *cache.Cache {
//
//	mutex.Lock()
//	stringOrderDataJson := string(orderDataJson)
//	err := cacheProgram.Add(orderId, stringOrderDataJson, cache.DefaultExpiration)
//	if err != nil {
//		println(err)
//	}
//
//	mutex.Unlock()
//
//	return cacheProgram
//
//}
