// Entrypoint for API
package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"shopping-cart/store"
)

func main() {
	port := "8002"
	router := store.NewRouter() // create routes
	// allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	// Launch server with CORS validations
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
