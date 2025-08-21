package controllers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	order, err := models.GetOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orderedItems, err := models.GetOrderedItems(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var total float32
	for _, item := range orderedItems {
		total += item.Price * float32(item.Quantity)
	}

	data := types.OrderPageData{
		Order:        *order,
		OrderedItems: orderedItems,
		Total:        total,
	}

	tmpl := template.Must(template.New("order").Funcs(template.FuncMap{
		"mul": func(a float32, b int) float32 { return a * float32(b) },
		"add": func(a float32, b float32) float32 { return a + b },
	}).ParseFiles(filepath.Join("pkg/views", "order.gohtml"),
		filepath.Join("pkg/views/components", "payment_card.gohtml")))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	buf.WriteTo(w)
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
	itemId, _ := strconv.Atoi(body.ItemID)
	oi := types.OrderedItem{
		OrderID: orderID,
		ItemID:  itemId,
	}
	er := models.DeleteOrderedItem(oi)
	if er != nil {
		http.Error(w, "Failed to delete ordered item", http.StatusInternalServerError)
		return
	}

	orderedItems, err := models.GetOrderedItems(orderID)
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

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Instructions string `json:"instructions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	order := types.Order{
		OrderID:             orderID,
		SpecificInstruction: body.Instructions,
		OrderStatus:         "In Progress",
	}

	err := models.UpdateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Order placed successfully"})
}
