package controllers

import(
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	items, err := models.GetAllItems()
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"categoryList":   types.CategoryList,
		"items":          items,
		"activeCategory": "All",
	}
	utils.WriteJSON(w, resp)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if !contains(types.CategoryList, category) {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	items, err := models.GetCategoryItems(category)
	if err != nil {
		http.Error(w, "Failed to fetch category items", http.StatusInternalServerError)
		return
	}

resp := map[string]interface{}{
		"categoryList":   types.CategoryList,
		"items":          items,
		"activeCategory": category,
	}
	utils.WriteJSON(w, resp)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		ItemID string `json:"item_id"`
		Qnty   string `json:"qnty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(body.Qnty)
	if err != nil || quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}
	orderID, exists := session.Values["order_id"].(int)
	if !exists {
		order := types.Order{
			TableNumber:      userID/50, //mp change this later
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
	}
	itemId,_:= strconv.Atoi(body.ItemID)

	existsInOrder, err := models.ItemExistsInOrder(itemId, orderID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	oi:=types.OrderedItem{
		OrderID: orderID,
		ItemID: itemId ,
		
	}

	if existsInOrder {
		err = models.UpdateOrderedItems(oi)
	} else {
		_,err = models.AddOrderedItem(oi)
	}
	if err != nil {
		http.Error(w, "Failed to add item", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Item added to cart"})

}


func CheckOrder(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	_, exists := session.Values["order_id"]
	utils.WriteJSON(w, map[string]bool{"exists": exists})
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}