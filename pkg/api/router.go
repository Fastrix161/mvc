package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/fastrix161/mvc/pkg/controllers"
	"github.com/fastrix161/mvc/pkg/middlewares"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Handle("/", http.RedirectHandler("/login", http.StatusFound))

	router.HandleFunc("/login", controllers.LoginHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/logout", controllers.Logout).Methods(http.MethodGet)
	router.HandleFunc("/signup", controllers.SignUpHandler).Methods(http.MethodGet, http.MethodPost)

	home := router.PathPrefix("/home").Subrouter()
	home.Use(middlewares.RestrictToLoggedIn)
	home.HandleFunc("/addToCart", controllers.AddToCart).Methods(http.MethodPost)  
	home.HandleFunc("/checkOrder", controllers.CheckOrder).Methods(http.MethodGet)
	home.HandleFunc("", controllers.GetHome).Methods(http.MethodGet)    
	home.HandleFunc("/{category}", controllers.GetCategory).Methods(http.MethodGet)  

	payment := router.PathPrefix("/payment").Subrouter()
	payment.Use(middlewares.RestrictToLoggedIn)
	payment.HandleFunc("", controllers.GetPayment).Methods(http.MethodGet)  
	payment.HandleFunc("/update", controllers.UpdatePayment).Methods(http.MethodPost) 

	chef := router.PathPrefix("/chef").Subrouter()
	chef.Use(middlewares.RestrictToChef)
	chef.HandleFunc("/orders", controllers.GetOrders).Methods(http.MethodGet)
	chef.HandleFunc("/change-status", controllers.ChangeStatus).Methods(http.MethodPatch)

	order := router.PathPrefix("/order").Subrouter()
	order.Use(middlewares.RestrictToLoggedIn)
	order.HandleFunc("", controllers.GetOrder).Methods(http.MethodGet)
	order.HandleFunc("/delete-item", controllers.DeleteOrderItem).Methods(http.MethodDelete)
	order.HandleFunc("/checkout", controllers.CheckoutOrder).Methods(http.MethodPost)

	return router
}

