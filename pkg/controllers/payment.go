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

func GetPayment(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	orderID, ok := session.Values["order_id"].(int)
	if !ok {
		http.Error(w, `{"message":"Order not found"}`, http.StatusNotFound)
		return
	}
	pay, err := models.GetPayment(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	validModes := map[string]bool{
		"Cash": true, "Card": true, "UPI": true, "Net Banking": true,
	}
	if pay == nil || !validModes[pay.Mode] {
		total, err := models.GetOrderTotal(orderID)
		if err != nil {
			http.Error(w, "Failed to calculate total", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "No payment created yet",
			"total":   total,
			"payment": nil,
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Please proceed to payment",
		"payment": pay,
	})
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)

	userID, ok := session.Values["user_id"].(int)
	if !ok || userID == 0 {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		PayID  string `json:"payment_id"`
		Mode   string `json:"mode"`
		Status bool   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !contains(types.Mode, body.Mode) || body.Mode == "unselected" {
		http.Error(w, "Invalid payment mode", http.StatusBadRequest)
		return
	}

	var pay types.Payment

	if body.PayID != "" {
		payID, _ := strconv.Atoi(body.PayID)
		payDB, err := models.GetPayment(payID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if payDB == nil {
			http.Error(w, "Payment not found", http.StatusNotFound)
			return
		}

		pay = types.Payment{
			PaymentID: payDB.PaymentID,
			OrderID:   payDB.OrderID,
			Total:     payDB.Total,
			Status:    body.Status,
			Mode:      body.Mode,
		}

		ok, err := models.UpdatePayment(pay)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Nothing to update"})
			return
		}

	} else {
		orderID, ok := session.Values["order_id"].(int)
		if !ok {
			http.Error(w, "Order ID not found in session", http.StatusBadRequest)
			return
		}

		total, err := models.GetOrderTotal(orderID)
		if err != nil {
			http.Error(w, "Failed to calculate order total", http.StatusInternalServerError)
			return
		}

		pay = types.Payment{
			OrderID: orderID,
			Total:   total,
			Mode:    body.Mode,
			Status:  body.Status,
		}

		payID, err := models.CreatePayment(pay)
		if err != nil {
			http.Error(w, "Failed to create payment", http.StatusInternalServerError)
			return
		}
		pay.PaymentID = payID
	}

	if body.Status {
		delete(session.Values, "order_id")
		session.Save(r, w)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Payment updated successfully"})
}
