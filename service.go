package main

import (
	"Taskl0/model"
	"database/sql"
)

type CacheService struct {
	cache Cache
	db    DB
}

func MakeCacheService() CacheService {
	return CacheService{
		cache: MakeCache(),
		db: DB{
			conn: new(sql.DB),
		},
	}
}

func (service *CacheService) AddOrder(model model.OrderModel) bool {
	order, _ := service.cache.GetOrderById(model.OrderUID)
	if order.OrderUID == "" {
		err := service.db.addOrderDataDB(model)
		if err != nil {
			return false
		}
		service.cache.AddOrder(model)
		return true
	}
	return false
}

func (service *CacheService) GetOrderById(idData string) (model.OrderModel, error) {
	return service.cache.GetOrderById(idData)
}

func (service *CacheService) Init() {
	orders, _ := service.db.getAllOrdersDataDB()
	for _, order := range orders {
		service.cache.AddOrder(order)
	}
}
