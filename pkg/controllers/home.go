package controllers

import (
	"bytes"
	"encoding/json"
	"html/template"
	_ "log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fastrix161/mvc/pkg/cache"
	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
	"github.com/gorilla/mux"
)

func GetJsonHome(w http.ResponseWriter, r *http.Request) {
	if cached, found := cache.Get("home_json"); found {
		utils.WriteJSON(w, http.StatusOK, cached)
		return
	}

	items, err := models.GetAllItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"categoryList":   types.CategoryList,
		"items":          items,
		"activeCategory": "All",
	}
	cache.Set("home_json", resp, 5*time.Minute)
	utils.WriteJSON(w, http.StatusOK, resp)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	if cached, found := cache.Get("home_html"); found {
		w.Write(cached.([]byte))
		return
	}
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
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, homepagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cache.Set("home_html", buf.Bytes(), 5*time.Minute)
	w.Write(buf.Bytes())
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

func GetSearchItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	items, err := models.GetItems(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	homepagedata := types.HomePage{
		ActiveCategory: "Search Results",
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
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Not logged in",
		})
		return
	}

	var body struct {
		ItemID string `json:"item_id"`
		Qnty   int    `json:"qnty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	quantity := body.Qnty
	if quantity <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Not logged in",
		})
		return
	}
	var orderID int
	orderIDInterface, exists := session.Values["order_id"]
	if !exists {
		order := types.Order{
			TableNumber:         rand.Intn(50) + 1,
			SpecificInstruction: "",
			OrderStatus:         "In Queue",
			UserID:              userID,
		}
		var err error
		orderID, err = models.AddOrder(order)
		if err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}
		session.Values["order_id"] = orderID
		session.Save(r, w)
	} else {
		orderID = orderIDInterface.(int)
	}

	itemId, err := strconv.Atoi(body.ItemID)
	if err != nil || itemId <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid item_id",
		})
		return
	}

	existsInOrder, err := models.ItemExistsInOrder(itemId, orderID)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}
	oi := types.OrderedItem{
		OrderID:  orderID,
		ItemID:   itemId,
		Quantity: quantity,
	}

	if existsInOrder {
		err = models.UpdateOrderedItems(oi)
	} else {
		_, err = models.AddOrderedItem(oi)
	}
	if err != nil {
		if err.Error() == "no changes to be made" {
			utils.WriteJSON(w, http.StatusOK, map[string]string{
				"message": "No changes to be made",
			})
			return
		}
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Item added to cart"})

}

func CheckOrder(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	_, exists := session.Values["order_id"]
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}
