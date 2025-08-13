package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"log"

	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
	"github.com/gorilla/sessions"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles(filepath.Join("pkg/views","login.gohtml")))
		err:= tmpl.Execute(w,nil)
		if err!=nil{
			http.Error(w, err.Error(),http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		
secret := os.Getenv("SECRET")
if secret == "" {
    log.Fatal("SECRET env variable not set")
}
var store = sessions.NewCookieStore([]byte(secret))
		user:=types.LoginUser{
			Email: r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		fmt.Println("Email:", user.Email, "Password:", user.Password)

		userDB, err := models.CheckEmail(user.Email)
		if err != nil || userDB.Email == "" {
			fmt.Println("User not found")		
			http.Redirect(w,r,"/login",http.StatusFound)
			return

		}
		check, err := utils.CheckPassword(user.Password, userDB.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !check {
			fmt.Println("Wrong password")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session, _ := store.Get(r, "session")
		session.Values["user_id"] = userDB.UserID
		err= session.Save(r, w)
		if err!=nil{
			fmt.Println("Failed to save session:", err)
        http.Error(w, "Error", http.StatusInternalServerError)
        return
		}
		payload := map[string]interface{}{
			"user_id": userDB.UserID,
			"email":   userDB.Email,
			"role":    userDB.Role,
		}
		token, err := utils.GenToken(payload)
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token_id",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(3 * time.Hour),
		})
		fmt.Println("User logged in successfully")
		switch userDB.Role {
		case "admin":
			http.Redirect(w, r, "/admin", http.StatusFound)
		case "chef":
			http.Redirect(w, r, "/chef", http.StatusFound)
		default:
			http.Redirect(w, r, "/home", http.StatusFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
