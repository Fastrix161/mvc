package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user types.SignupUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Errorf("invalid json format: %v", err).Error(), http.StatusBadRequest)
		return
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

	// usercheck, _ := models.CheckEmail(userDB.Email)
	// if usercheck != nil {
	// 	http.Error(w,)
	// }
	_, er := models.CreateUser(userDB)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered Successfully"))

}
