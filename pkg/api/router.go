package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/fastrix161/mvc/pkg/controllers"
	"github.com/fastrix161/mvc/pkg/middlewares"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/login", http.StatusFound))
	return router
}

func HomeRoutes(mux *http.ServeMux) {
	mux.Handle("/home", middlewares.RestrictToLoggedIn(methodHandler("GET", controllers.GetHome)))
	mux.Handle("/home/", middlewares.RestrictToLoggedIn(methodHandler("GET", controllers.GetCategory)))
	mux.Handle("/home/addToCart", methodHandler("POST", controllers.AddToCart))
	mux.Handle("/home/checkOrder", middlewares.RestrictToLoggedIn(methodHandler("GET", controllers.CheckOrder)))
}

func LoginRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/logout", controllers.Logout)
}

func LogoutRoutes(mux *http.ServeMux) {
	mux.Handle("/logout", methodHandler(http.MethodPost, controllers.LoginHandler))
}

func SignupRoute(mux *http.ServeMux){
	mux.Handle("/signup", methodHandler(http.MethodPost, controllers.SignUp))
}

func PaymentRoutes(mux *http.ServeMux) {
mux.Handle("/payment", middlewares.RestrictToLoggedIn(http.HandlerFunc(controllers.GetPayment)))
mux.Handle("/payment/update", middlewares.RestrictToLoggedIn(http.HandlerFunc(controllers.UpdatePayment)))
}
func ChefRouter(mux *http.ServeMux) {
	mux.Handle("/chef/orders", middlewares.RestrictToChef(http.HandlerFunc(controllers.GetOrders)))
	mux.Handle("/chef/change-status", middlewares.RestrictToChef(http.HandlerFunc(controllers.ChangeStatus)))
}

func RegisterOrderRoutes(mux *http.ServeMux) {
	mux.Handle("/order", middlewares.RestrictToLoggedIn(http.HandlerFunc(controllers.GetOrder)))

	mux.Handle("/order/delete-item", middlewares.RestrictToLoggedIn(http.HandlerFunc(controllers.DeleteOrderItem)))

	mux.Handle("/order/checkout", middlewares.RestrictToLoggedIn(http.HandlerFunc(controllers.CheckoutOrder)))
}

func methodHandler(method string, h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}