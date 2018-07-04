package store

import (
	"net/http"
	"testing"
	//"fmt"
	"gopkg.in/gavv/httpexpect.v1"
)

// Test cart addition
func TestCartAddGetCalculateDeleteInventory(t *testing.T) {
	e := httpexpect.New(t, API_URL)
	// Shopping Cart
	product := map[string]map[string]interface{}{
		"cart": {
			"Belts":    3,
			"Trousers": 2,
			"Shoes":    1,
			"Shirts":   3,
			"Ties":     8,
		},
	}
	// StatusUnauthorized check without valid token
	e.POST("/shopping/cart/AddUserCart/1235").WithHeader("Authorization", "Bearer "+"dummy").
		WithJSON(product).
		Expect().
		Status(http.StatusUnauthorized)

	// get token
	credentials := map[string]interface{}{
		"username": "user2",
		"password": "password2",
	}
	r := e.POST("/get-token").
		WithJSON(credentials).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token := r.Value("token").String().Raw()
	// add product
	e.POST("/shopping/cart/AddUserCart/1235").WithHeader("Authorization", "Bearer "+token).
		WithJSON(product).
		Expect().
		Status(http.StatusCreated)
	// calculate cart
	e.GET("/shopping/cart/calculate/1235").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).Body().Equal("538")
	uproduct := map[string]interface{}{
		"ID": 1235,
		"cart": map[string]interface{}{
			"Belts":    3,
			"Trousers": 2,
			"Shoes":    1,
			"Shirts":   2,
			"Ties":     8,
		},
	}
	// update cart
	e.PUT("/shopping/cart/UpdateUserCart/1235").WithHeader("Authorization", "Bearer "+token).
		WithJSON(uproduct).
		Expect().
		Status(http.StatusOK)
	// calculate cart
	e.GET("/shopping/cart/calculate/1235").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).Body().Equal("573")
	// delete cart
	e.DELETE("/shopping/cart/DeleteUserCart/1235").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)
}
