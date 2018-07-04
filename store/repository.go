package store

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "dummyStore"

// COLLECTION is the name of the collection in DB
const COLLECTION = "store"
const CARTCOLLECTION = "cart"
const ADMINCOLLECTION = "adminCart"
const USERCOLLECTION = "users"

var productId = 6

// GetAdmin validates user
func (r Repository) GetAdmin(userName string, password string) bool {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(ADMINCOLLECTION)
	if count, err := c.Find(bson.M{"username": userName, "password": password}).Count(); err != nil || count == 0 {
		fmt.Println("Failed to write results:", err)
		return false
	}
	return true
}

// GetProducts returns the list of Products
func (r Repository) GetProducts() Products {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Products{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// Get products in Map Structure
func (r Repository) GetMapProducts() MapProduct {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Products{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}
	mresults := MapProduct{}
	for _, product := range results {
		mresults[product.ItemName] = product
	}
	return mresults
}

// GetProductById returns a unique Product
func (r Repository) GetProductById(id int) Product {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result Product

	fmt.Println("ID in GetProductById", id)

	if err := c.FindId(id).One(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// GetProductsByString takes a search string as input and returns products
func (r Repository) GetProductsByString(query string) Products {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	result := Products{}

	// Logic to create filter
	qs := strings.Split(query, " ")
	and := make([]bson.M, len(qs))
	for i, q := range qs {
		and[i] = bson.M{"itemname": bson.M{
			"$regex": bson.RegEx{Pattern: ".*" + q + ".*", Options: "i"},
		}}
	}
	filter := bson.M{"$and": and}

	if err := c.Find(&filter).Limit(5).All(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// AddProduct adds a Product in the DB
func (r Repository) AddProduct(product Product) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	productId += 1
	product.ID = productId
	session.DB(DBNAME).C(COLLECTION).Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Product ID- ", product.ID)

	return true
}

// UpdateProduct updates a Product in the DB
func (r Repository) UpdateProduct(product Product) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	err = session.DB(DBNAME).C(COLLECTION).UpdateId(product.ID, product)

	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Updated Product ID - ", product.ID)

	return true
}

// DeleteProduct deletes an Product
func (r Repository) DeleteProduct(id int) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// Remove product
	if err = session.DB(DBNAME).C(COLLECTION).RemoveId(id); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	fmt.Println("Deleted Product ID - ", id)
	// Write status
	return "OK"
}

// GetUser validates user
func (r Repository) GetUser(userName string, password string) bool {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(USERCOLLECTION)

	if count, err := c.Find(bson.M{"username": userName, "password": password}).Count(); err != nil || count == 0 {
		fmt.Println("Failed to write results:", err)
		return false
	}
	return true
}

// GetUser validates user
func (r Repository) GetUserById(id int, userName string) bool {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(USERCOLLECTION)
	if count, err := c.Find(bson.M{"_id": id, "username": userName}).Count(); err != nil || count == 0 {
		fmt.Println("Failed to write results:", err)
		return false
	}
	return true
}

// GetUserCartById returns a cart specific to the user
func (r Repository) GetUserCartById(id int) Cart {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(CARTCOLLECTION)
	var result Cart

	if err := c.FindId(id).One(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// AddCart adds a cart in the DB
func (r Repository) AddUserCart(cart Cart, id int) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	cart.ID = id
	session.DB(DBNAME).C(CARTCOLLECTION).Insert(cart)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New cart ID- ", cart.ID)

	return true
}

// DeleteProduct deletes an Product
func (r Repository) DeleteUserCart(id int) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// Remove product
	if err = session.DB(DBNAME).C(CARTCOLLECTION).RemoveId(id); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	fmt.Println("Deleted UserCart ID - ", id)
	// Write status
	return "OK"
}

// UpdateUserCart updates a cart in the DB
func (r Repository) UpdateUserCart(cart Cart) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	err = session.DB(DBNAME).C(CARTCOLLECTION).UpdateId(cart.ID, cart)

	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Updated cart ID - ", cart.ID)

	return true
}
