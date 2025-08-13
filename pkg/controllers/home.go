package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
	"github.com/gorilla/mux"
)

// func GetHome(w http.ResponseWriter, r *http.Request) {

// 	items, err := models.GetAllItems()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp := map[string]interface{}{
// 		"categoryList":   types.CategoryList,
// 		"items":          items,
// 		"activeCategory": "All",
// 	}
// 	utils.WriteJSON(w, resp)
// }

func GetHome(w http.ResponseWriter, r *http.Request) {

	items, err := models.GetAllItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	homepagedata := types.HomePage{
		ActiveCategory: "All",
		CategoryList:   types.CategoryList,
		Items:          items,
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "home.gohtml")))
	err = tmpl.Execute(w, homepagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	if !contains(types.CategoryList, category) {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	items, err := models.GetCategoryItems(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	homepagedata := types.HomePage{
		ActiveCategory: category,
		CategoryList:   types.CategoryList,
		Items:          items,
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "home.gohtml")))
	err = tmpl.Execute(w, homepagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	userID, ok := session.Values["user_id"].(int)
	if !ok {
	utils.WriteJSON(w, map[string]string{
		"error": "Not logged in",
	})
	w.WriteHeader(http.StatusUnauthorized)
	return
}

	var body struct {
		ItemID string `json:"item_id"`
		Qnty   int `json:"qnty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {

	utils.WriteJSON(w, map[string]string{
		"error": err.Error(),
	})
	w.WriteHeader(http.StatusUnauthorized)
	return
}

	quantity :=body.Qnty
	if quantity <= 0 {
	utils.WriteJSON(w, map[string]string{
		"error": "Not logged in",
	})
	w.WriteHeader(http.StatusUnauthorized)
	return
}
	orderID, exists := session.Values["order_id"].(int)
	_ = session.Save(r, w)
	if !exists {
		order := types.Order{
			TableNumber:         userID % 50, //mp change this later
			SpecificInstruction: "",
			OrderStatus:         "In Queue",
			UserID:              userID,
		}
		orderID, err := models.AddOrder(order)
		if err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}
		session.Values["order_id"] = orderID
		session.Save(r,w)
	}
	itemId, _ := strconv.Atoi(body.ItemID)

	existsInOrder, err := models.ItemExistsInOrder(itemId, orderID)
	if err != nil {
		utils.WriteJSON(w, map[string]string{
		"error": err.Error(),
	})
	w.WriteHeader(http.StatusUnauthorized)
	return
}
	oi := types.OrderedItem{
		OrderID: orderID,
		ItemID:  itemId,
		Quantity: quantity,
	}

	if existsInOrder {
		err = models.UpdateOrderedItems(oi)
	} else {
		_, err = models.AddOrderedItem(oi)
	}
	if err != nil {
		utils.WriteJSON(w, map[string]string{
		"error": err.Error(),
	})
	w.WriteHeader(http.StatusUnauthorized)
	return
	}
	utils.WriteJSON(w, map[string]string{"message": "Item added to cart"})
	 
}

func CheckOrder(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	_, exists := session.Values["order_id"]
	if !exists{
		utils.WriteJSON(w, map[string]string{"message":"order not created"})
	}
	utils.WriteJSON(w, map[string]bool{"exists": exists})
}
