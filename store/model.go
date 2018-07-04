package store

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminUser struct {
	ID       int    `bson:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type NormalUser struct {
	ID       int    `bson:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type Product struct {
	ID       int    `bson:"_id"`
	ItemName string `json:"itemname"`
	Stock    int    `json:"stock"`
	Price    int    `json:"price"`
}

// Products is an array of Product objects
type Products []Product
type MapProduct map[string]Product

// Shopping Cart
type Cart struct {
	ID   int            `bson:"_id"`
	Cart map[string]int `json:"cart"`
}

// Offers Cart
type Offer struct {
	ID    int                          `bson:"_id"`
	Offer map[string]map[string]string `json:"offer"`
}

type Offers []Offer
