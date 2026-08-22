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

func GetAdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "admin.gohtml")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := models.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	userCall, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !contains(types.Role, body.Role) {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(body.UserID)
	if err != nil || userId <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user_id",
		})
		return
	}

	if userId == userCall {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Security: Can not change Role of self!!",
		})
		return
	}

	err = models.SetUserRole(userId, body.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User role updated"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session := middlewares.GetSession(r)
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	userId, err := strconv.Atoi(body.UserID)
	if err != nil || userId <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
		return
	}
	user, err := models.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user.Role == "admin" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Admin can't be deleted"})
		return
	}

	err = models.DeleteUser(userId)
	if err != nil {
		log.Println("TRUE DELETION ERROR:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
