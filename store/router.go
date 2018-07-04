package store

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"AuthenticationAdmin",
		"POST",
		"/get-admin-token",
		controller.GetAdminToken,
	},
	Route{
		"AuthenticationUser",
		"POST",
		"/get-token",
		controller.GetToken,
	},
	Route{
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	Route{
		"AddProduct",
		"POST",
		"/AddProduct",
		AuthenticationMiddleware(controller.AddProduct),
	},
	Route{
		"UpdateProduct",
		"PUT",
		"/UpdateProduct",
		AuthenticationMiddleware(controller.UpdateProduct),
	},
	// Get Product by {id}
	Route{
		"GetProduct",
		"GET",
		"/products/{id}",
		controller.GetProduct,
	},
	// Delete Product by {id}
	Route{
		"DeleteProduct",
		"DELETE",
		"/deleteProduct/{id}",
		AuthenticationMiddleware(controller.DeleteProduct),
	},
	// Search product with string
	Route{
		"SearchProduct",
		"GET",
		"/Search/{query}",
		controller.SearchProduct,
	},
	// User CartView
	Route{
		"GetUserCart",
		"GET",
		"/shopping/cart/{id}",
		controller.GetUserCart,
	},
	// User Add Cart
	Route{
		"AddUserCart",
		"POST",
		"/shopping/cart/AddUserCart/{id}",
		AuthenticationMiddlewareUser(controller.AddUserCart),
	},
	// Update User Cart
	Route{
		"UpdateUserCart",
		"PUT",
		"/shopping/cart/UpdateUserCart/{id}",
		AuthenticationMiddlewareUser(controller.UpdateUserCart),
	},
	// User Delete Cart
	Route{
		"DeleteUserCart",
		"DELETE",
		"/shopping/cart/DeleteUserCart/{id}",
		AuthenticationMiddlewareUser(controller.DeleteUserCart),
	},
	//Calculate User cart
	Route{
		"CalculateUserCart",
		"GET",
		"/shopping/cart/calculate/{id}",
		controller.CalculateUserCart,
	}}

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
