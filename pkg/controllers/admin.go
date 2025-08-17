 package controllers

import (
	// "encoding/json"
	"net/http"
	// "strconv"
	"path/filepath"
	"html/template"

	
	// "github.com/fastrix161/mvc/pkg/middlewares"
	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	// "github.com/fastrix161/mvc/pkg/utils"
)

func GetAdminPage(w http.ResponseWriter, r *http.Request){
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orders, err := models.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminpagedata:=types.AdminPage{
		Users: users,
		Orders: orders,
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "admin.gohtml")))
	err = tmpl.Execute(w, adminpagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllOrders(w http.ResponseWriter, r *http.Request){
	orders, err := models.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminpagedata:=types.AdminOrderPage{
		Orders: orders,
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "admin.gohtml")))
	err = tmpl.Execute(w, adminpagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request){
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminpagedata:=types.AdminUserPage{
		Users: users,
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "admin.gohtml")))
	err = tmpl.Execute(w, adminpagedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

