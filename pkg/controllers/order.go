package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func GetOrder(w http.ResponseWriter, r *http.Request){
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Error(w, `{"message":"Order not found"}`, http.StatusNotFound)
		return
	}
	order, err := models.GetOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}
	orderedItems, err := models.GetOrderedItems(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, map[string]interface{}{
		"order":        order,
		"orderedItems": orderedItems,
	})
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Error(w, `{"message":"Order not found"}`, http.StatusNotFound)
		return
	}

	var body struct {
		ItemID string `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	itemId,_:= strconv.Atoi(body.ItemID)
oi:= types.OrderedItem{
	OrderID :orderID,
	ItemID: itemId ,
}
	er := models.DeleteOrderedItem(oi)
	if er != nil {
		http.Error(w, "Failed to delete ordered item", http.StatusInternalServerError)
		return
	}

	orderedItems, err :=  models.GetOrderedItems(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
		return
	}

	if len(orderedItems) == 0 {
		err := models.DeleteOrder(orderID)
		if err != nil {
			http.Error(w, "Failed to delete order", http.StatusInternalServerError)
			return
		}
		delete(session.Values, "order_id")
		session.Save(r, w)
		utils.WriteJSON(w, map[string]string{"message": "Cart is empty, order deleted"})
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Item deleted"})
}

func CheckoutOrder(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Error(w, `{"message":"Order not found"}`, http.StatusNotFound)
		return
	}

	order, err := models.GetOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
		return
	}
	orderedItems, err := models.GetOrderedItems(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
		return
	}
	
	utils.WriteJSON(w, map[string]interface{}{
		"order":        order,
		"orderedItems": orderedItems,
	})
}
