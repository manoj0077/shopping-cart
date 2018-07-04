package store

import (
	"net/http"
	"testing"
	//"fmt"
	"gopkg.in/gavv/httpexpect.v1"
)

const (
	API_URL = `http://localhost:8002`
)

//Test Basic Get Products Inventory
func TestProductsInventory(t *testing.T) {
	e := httpexpect.New(t, API_URL)
	schema := `{
				"type": "array",
				"items": {
						"type": "object",
						"properties": {
								"ID": {
										"type": "integer"
								},
								"itemname": {
										"type": "string"
								},
								"stock": {
										"type": "integer"
								},
								"price": {
										"type": "integer"
								}
						}
				}
		}`
	products := e.GET("/").
		Expect().
		Status(http.StatusOK).JSON()
	products.Schema(schema)
}

//Test Admin Token validation and add inventory
func TestAdminAddSearchDeleteInventory(t *testing.T) {
	product := map[string]interface{}{
		"itemname": "Test",
		"stock":    10,
		"price":    20,
	}
	e := httpexpect.New(t, API_URL)
	// StatusUnauthorized check without valid token
	e.POST("/AddProduct").WithHeader("Authorization", "Bearer "+"dummy").
		WithJSON(product).
		Expect().
		Status(http.StatusUnauthorized)
		// get token
	credentials := map[string]interface{}{
		"username": "admin",
		"password": "password",
	}
	r := e.POST("/get-admin-token").
		WithJSON(credentials).
		Expect().
		Status(http.StatusOK).JSON().Object()
	token := r.Value("token").String().Raw()
	// add product
	e.POST("/AddProduct").WithHeader("Authorization", "Bearer "+token).
		WithJSON(product).
		Expect().
		Status(http.StatusCreated)
	// get product /Search/{query}
	sproduct := e.GET("/Search/Test").
		Expect().
		Status(http.StatusOK).JSON().Array()
	id := sproduct.Element(0).Object().Value("ID").Raw()
	e.DELETE("/deleteProduct/{id}", id.(float64)).WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)
}
