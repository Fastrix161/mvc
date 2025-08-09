package middlewares

import(
	"net/http"
	"fmt"

	
	"github.com/fastrix161/mvc/pkg/utils"
)

func Log(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) (string,	error){
	cookie, err:= r.Cookie("token_id")
	if err !=nil{
		return "",fmt.Errorf("error getting cookie: %v",err)
	}
	return cookie.Value, nil
}

func RestrictToAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		token,err:= getToken(r)
		if err!= nil{
			http.Redirect(w,r,"/login", http.StatusFound)
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		role,_:=claims["role"].(string)
		if role !="admin"{
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		
		ctx := setUserInContext(r, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})}

	func RestrictToChef(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil || token == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		role, _ := claims["role"].(string)
		if role != "chef" && role != "admin" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := setUserInContext(r, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

	func RestrictToLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		token,err:= getToken(r)
		if err!= nil{
			http.Redirect(w,r,"/login", http.StatusFound)
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		
		ctx := setUserInContext(r, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})}

	func RestrictToNew(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err == nil && token != "" {
			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}