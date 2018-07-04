package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Controller ...
type Controller struct {
	Repository Repository
}

/* Middleware handler to handle all admin requests for authentication */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret123"), nil
				})
				if error != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					log.Println("TOKEN WAS VALID")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

// Get AdminAuthentication token /
func (c *Controller) GetAdminToken(w http.ResponseWriter, req *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576)) // read the body of the request

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := req.Body.Close(); err != nil {
		log.Fatalln("Error Get Token", err)
	}

	if err := json.Unmarshal(body, &user); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			return
		}
	}

	// Validate AdminUser
	valid := c.Repository.GetAdmin(user.Username, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: "Not Valid Admin Credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	tokenString, error := token.SignedString([]byte("secret123"))
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// Get Authentication token GET /
func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576)) // read the body of the request

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := req.Body.Close(); err != nil {
		log.Fatalln("Error Get Token", err)
	}

	if err := json.Unmarshal(body, &user); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			return
		}
	}
	// Validate User
	valid := c.Repository.GetUser(user.Username, user.Password)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: "Not Valid Credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

/* Middleware handler to handle all requests for authentication */
func AuthenticationMiddlewareUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				claims, ok := token.Claims.(jwt.MapClaims)
				if !(ok || token.Valid) {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
				var user User
				user.Username = claims["username"].(string)
				user.Password = claims["password"].(string)
				vars := mux.Vars(req)
				id := vars["id"]
				userid, err := strconv.Atoi(id)
				if err != nil {
					log.Fatalln("Error GetUserId", err)
				}
				var controller = &Controller{Repository: Repository{}}
				valid := controller.Repository.GetUserById(userid, user.Username)
				if token.Valid && valid {
					//log.Println("TOKEN WAS VALID")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	products := c.Repository.GetProducts() // list of all products
	// log.Println(products)
	data, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddProduct POST /
func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	//log.Println(body)

	if err != nil {
		log.Fatalln("Error AddProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddProduct", err)
	}

	if err := json.Unmarshal(body, &product); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			return
		}
	}

	success := c.Repository.AddProduct(product) // adds the product to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// SearchProduct GET /
func (c *Controller) SearchProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	query := vars["query"] // param query

	products := c.Repository.GetProductsByString(query)
	data, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdateProduct PUT /
func (c *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error UpdateProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdateProduct", err)
	}

	if err := json.Unmarshal(body, &product); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateProduct unmarshalling data", err)
			return
		}
	}
	success := c.Repository.UpdateProduct(product) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// GetProduct GET - Gets a single product by ID /
func (c *Controller) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"] // param id

	productid, err := strconv.Atoi(id)

	if err != nil {
		log.Fatalln("Error GetProduct", err)
	}

	product := c.Repository.GetProductById(productid)
	data, _ := json.Marshal(product)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// DeleteProduct DELETE /
func (c *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] // param id

	productid, err := strconv.Atoi(id)

	if err != nil {
		log.Fatalln("Error GetProduct", err)
	}

	if err := c.Repository.DeleteProduct(productid); err != "" { // delete a product by id
		log.Println(err)
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// Gets a cart by userID /
func (c *Controller) GetUserCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"] // param id

	userid, err := strconv.Atoi(id)

	if err != nil {
		log.Fatalln("Error GetUserCart", err)
	}

	cart := c.Repository.GetUserCartById(userid)
	data, _ := json.Marshal(cart)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// Modify Cart /
func (c *Controller) AddUserCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] // param id

	userid, err := strconv.Atoi(id)

	if err != nil {
		log.Fatalln("Error userID", err)
	}

	var cart Cart
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	if err != nil {
		log.Fatalln("Error AddUserCart", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddUserCart", err)
	}

	if err := json.Unmarshal(body, &cart); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddUserCart unmarshalling data", err)
			return
		}
	}

	// check whether cart is in limits or not
	valid := c.ValidCart(cart)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Added more products than inventory"))
		return
	}

	success := c.Repository.AddUserCart(cart, userid) // adds the product to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// DeleteUserCart DELETE /
func (c *Controller) DeleteUserCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] // param id

	userid, err := strconv.Atoi(id)

	if err != nil {
		log.Fatalln("Error GetUserId", err)
	}

	if err := c.Repository.DeleteUserCart(userid); err != "" { // delete a cart by id
		log.Println(err)
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// UpdateProduct PUT /
func (c *Controller) UpdateUserCart(w http.ResponseWriter, r *http.Request) {
	var cart Cart
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error UpdateUserCart", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdateUserCart", err)
	}
	if err := json.Unmarshal(body, &cart); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateProduct unmarshalling data", err)
			return
		}
	}

	// check whether cart is in limits or not
	valid := c.ValidCart(cart)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Added more products than inventory"))
		return
	}

	success := c.Repository.UpdateUserCart(cart) // updates the cart in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// Calculate cart by applying promotions
func (c *Controller) CalculateUserCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] // param id

	userid, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln("Error GetUserId", err)
	}

	costCart := c.PopulateCart(userid)

	offers := getOffers()

	for i := range offers {
		costCart = offers[i].calculateDiscountPrice(costCart)
		//fmt.Println(costCart)
	}
	amount := c.CartValue(costCart)
	data, _ := json.Marshal(amount)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (c *Controller) ValidCart(cart Cart) bool {
	products := c.Repository.GetMapProducts()
	for k, v := range cart.Cart {
		prod, ok := products[k]
		if !ok || prod.Stock < v {
			return false
		}
	}
	return true
}

func (c *Controller) PopulateCart(userid int) CostCart {
	costCart := make(map[string]AmountCart)
	cart := c.Repository.GetUserCartById(userid)
	products := c.Repository.GetMapProducts()
	for k, v := range cart.Cart {
		prodCart := AmountCart{Quantity: v, Price: products[k].Price, Amount: float64(v * products[k].Price)}
		costCart[k] = prodCart
	}
	return costCart
}

func (c *Controller) CartValue(cart CostCart) float64 {
	var amount float64
	for k, _ := range cart {
		amount += cart[k].Amount
	}
	return amount
}
