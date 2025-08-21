package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)
func GetOrders(w http.ResponseWriter, r *http.Request) {
    orders, err := models.GetAllOrders()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var ordersWithItems []types.OrderWithItems
    for _, order := range orders {
        items, err := models.GetOrderedItems(order.OrderID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ordersWithItems = append(ordersWithItems, types.OrderWithItems{
            Order:        order,
            OrderedItems: items,
        })
    }

    chefpagedata := types.ChefPage{
        Orders: ordersWithItems,
    }

    tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "chef.gohtml")))
    err = tmpl.Execute(w, chefpagedata)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		OrderID string `json:"order_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !contains(types.OrderStatus, body.Status) {
		http.Error(w, "Invalid status, only \"In Queue\", \"In Progress\" and \"Completed\" allowed", http.StatusBadRequest)
		return
	}

	orderId,_:= strconv.Atoi(body.OrderID)
	orderDB,err := models.GetOrder(orderId)
	if err!=nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	order:= types.Order{
		OrderID: orderDB.OrderID,
		TableNumber: orderDB.TableNumber,
		OrderStatus: body.Status,
		UserID: orderDB.UserID,
		SpecificInstruction: orderDB.SpecificInstruction,
	}

	er:= models.UpdateOrder(order)
	if er!=nil{
		if er.Error()=="order not found"{
			utils.WriteJSON(w,map[string]string{"message": "No change to make"} )
			return
		}
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w,map[string]string{"message": "Order status updated successfully"})
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}