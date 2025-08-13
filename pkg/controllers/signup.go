package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "signup.gohtml")))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case http.MethodPost:
		
	mobnum, err := strconv.Atoi(r.FormValue("mobile_number"))
	if err!=nil{
		http.Error(w,"Mobile number should only contain digits", http.StatusBadRequest)
		return
	}
		user:=types.SignupUser{
			Email: r.FormValue("email"),
			Password: r.FormValue("password") ,
			MobileNumber: mobnum,
			Name:  r.FormValue("name"),
		}
		if user.Name == "" || user.Email == "" || user.Password == "" {
			http.Error(w, fmt.Errorf("name, email and password can't be empty").Error(), http.StatusBadRequest)
			return
		}
		if len(user.Password) < 8 || len(user.Password) > 20 {
			http.Error(w, fmt.Errorf("password length should be between 8 and 20").Error(), http.StatusBadRequest)
			return
		}
		name := user.Name
		mobile_number := user.MobileNumber
		email := user.Email
		password := user.Password

		hashed_pwd, err := utils.GenHash(password, 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userDB := types.User{
			Name:         name,
			MobileNumber: mobile_number,
			Role:         "customer",
			Email:        email,
			Password:     hashed_pwd,
		}

		usercheck, _ := models.CheckEmail(userDB.Email)
		if usercheck != nil {
			http.Error(w,"User Exists already",http.StatusBadRequest)
			return
		}
		_, er := models.CreateUser(userDB)
		if er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User registered Successfully"))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	}
}
