package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func GetPayment(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Error(w, `{"message":"Order not found"}`, http.StatusNotFound)
		return
	}
	pay,err:=models.GetPayment(orderID)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pay==nil{
		total, err := models.GetOrderTotal(orderID)
		if err != nil {
			http.Error(w, "Failed to calculate total", http.StatusInternalServerError)
			return
		}

		newPayment := types.Payment{
			OrderID: orderID,
			Total:   total,
			Mode:    "unselected", 
			Status:  false, 
		}

		insertID, err := models.CreatePayment(newPayment)
		if err != nil {
			http.Error(w, "Failed to create payment", http.StatusInternalServerError)
			return
		}

		newPayment.PaymentID = insertID
		pay = &newPayment
	}
	
	utils.WriteJSON(w, map[string]interface{}{
		"message": "Please proceed to payment",
		"payment": pay,
	})
}

func UpdatePayment(w http.ResponseWriter, r *http.Request){
	session := middlewares.GetSession(r)
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		PayID string `json:"payment_id"`
		Mode string `json:"mode"`
		Status bool `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	
	// if !contains(types.Mode, body.Mode) {
	// 	http.Error(w, "Invalid mode, only \"unselected\",\"Cash\",\"Card\",\"UPI\",\"Net Banking\" allowed", http.StatusBadRequest)
	// 	return
	// }
	
	payId,_:= strconv.Atoi(body.PayID)
	payDB,err := models.GetPayment(payId)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	pay:= types.Payment{
		PaymentID: payDB.PaymentID,
		OrderID: payDB.OrderID,
		Total: payDB.Total,
		Status: body.Status,
		Mode: body.Mode,
	}

	check,er:= models.UpdatePayment(pay)
	if er!=nil{
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	if !check{
		utils.WriteJSON(w,map[string]string{"message": "Nothing to Update"})
	}
	utils.WriteJSON(w,map[string]string{"message": "Payment updated"})
}