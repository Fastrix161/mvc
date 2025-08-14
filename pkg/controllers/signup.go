package controllers

import (
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
		renderSignupPage(w,"")
		return
	case http.MethodPost:
		
	mobnum, err := strconv.Atoi(r.FormValue("mobile_number"))
	if err!=nil{
		renderSignupPage(w, "Mobile number should only contain digits")
		return
	}
		user:=types.SignupUser{
			Email: r.FormValue("email"),
			Password: r.FormValue("password") ,
			MobileNumber: mobnum,
			Name:  r.FormValue("name"),
		}
		if user.Name == "" || user.Email == "" || user.Password == "" {
			renderSignupPage(w, "Name, email, and password can't be empty")
			return
		}
		if len(user.Password) < 8 || len(user.Password) > 20 {
			renderSignupPage(w, "Password must be between 8 and 20 characters")
			return
		}
		name := user.Name
		mobile_number := user.MobileNumber
		email := user.Email
		password := user.Password

		hashed_pwd, err := utils.GenHash(password, 10)
		if err != nil {
			renderSignupPage(w, "Something went wrong. Please try again.")
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
			renderSignupPage(w, "User already exists with this email.")
			return
		}
		_, er := models.CreateUser(userDB)
		if er != nil {
			renderSignupPage(w, "Failed to create user.")
			return
		}
		http.Redirect(w, r, "/login?success=1", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	}
}
func renderSignupPage(w http.ResponseWriter, errorMessage string) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views", "signup.gohtml")))
	err := tmpl.Execute(w, struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}