package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type ViewData struct {
	OrderId   string
	OrderInfo string
}

func serverHtmlStart(service *CacheService) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handled")

		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handled /postform")

		orderId := r.FormValue("getOrder")

		model, flag := service.GetOrderById(orderId)

		var data ViewData

		if flag == nil {

			var strOrderJson string = fmt.Sprintf("%+v", model)

			var orderMapJson map[string]interface{}
			json.Unmarshal([]byte(strOrderJson), &orderMapJson)

			var orderId string = fmt.Sprintf("%+v", orderMapJson["order_uid"])

			data = ViewData{
				OrderId:   orderId,
				OrderInfo: strOrderJson,
			}

		} else {
			data = ViewData{
				OrderId:   "нет информации",
				OrderInfo: "нет информации",
			}
		}

		tmpl, _ := template.ParseFiles("web/getOrder.html")

		tmpl.Execute(w, data)

	})

	http.ListenAndServe(":9090", nil)
}
