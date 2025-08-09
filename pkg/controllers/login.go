package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"

	"github.com/fastrix161/mvc/pkg/models"
	"github.com/fastrix161/mvc/pkg/types"
	"github.com/fastrix161/mvc/pkg/utils"
	"github.com/gorilla/sessions"
)


var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET")))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case http.MethodGet:
	// 	RenderTemplate(w, "login.html", nil)
	// 	return

	case http.MethodPost:

		var user types.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
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
		}
		if !check {
			fmt.Println("Wrong password")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session, _ := store.Get(r, "session-name")
		session.Values["user_id"] = userDB.UserID
		session.Save(r, w)

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

	}
}
