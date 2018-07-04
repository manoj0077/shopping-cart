# shopping-cart
i) Download or clone repository into your $GOPATH/src/ 

ii) Install MongoDB and start the service

iii) Open the terminal, run below commands to import data for inventory, admin, users

	cd <$GOPATH/src/shopping-cart>
	mongo < dummyInventory.js
	mongo < dummyAdmin.js
	mongo < dummyUsers.js
	
iv) Install go dependencies by running below commands
	
	go get "github.com/dgrijalva/jwt-go"
	go get "github.com/gorilla/mux"
	go get "github.com/gorilla/context"
	go get "github.com/gorilla/handlers"
	go get "gopkg.in/mgo.v2"
	go get gopkg.in/gavv/httpexpect.v1
	
v) Now run the API using 

	go run main.go
	
vi) To Run test cases, 

	cd <$GOPATH/src/shopping-cart/test-cases
	go test -v
	
---------------------------------

API for admin and product inventory:- 

i) Get product inventory:- 

	GET http://localhost:8002

	sample output:-

	[
		{
			"ID": 1,
			"itemname": "Belts",
			"stock": 10,
			"price": 20
		},
		{
			"ID": 2,
			"itemname": "Shirts",
			"stock": 5,
			"price": 60
		}
	]

ii) Get Admin Token to add products, update or delete products 

	POST http://localhost:8002/get-admin-token

	BODY:- 
	{ "username": "admin", "password": "password"}

	sample output:-
	{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6InBhc3N3b3JkIiwidXNlcm5hbWUiOiJhZG1pbiJ9.wyj7wX6Vit8RtYblcwcX_5shG5PQ90FmpNEe7fnBOv4"
	}

iii) Add product to inventory

	POST http://localhost:8002/AddProduct

	Headers:-
	Authorization Bearer <tokenValue>

	BODY:-
	{
		"itemname": "Test",
		"stock": 10,
		"price": 20
	}

iv) Delete Product 

	DELETE http://localhost:8002/deleteProduct/7

	Headers:-
	Authorization Bearer <tokenValue>

------

API for users and addition, view, calculate, update, delete cart

i) Get User Token to add, update or delete cart

	POST http://localhost:8002/get-token

	BODY:- 
	{ "username": "user1", "password": "password1"}

	sample output:-
	{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6InBhc3N3b3JkIiwidXNlcm5hbWUiOiJhZG1pbiJ9.wyj7wX6Vit8RtYblcwcX_5shG5PQ90FmpNEe7fnBOv4"
	}


ii) Add Cart

	POST http://localhost:8002/shopping/cart/AddUserCart/1234

	Headers:-
	Authorization Bearer <tokenValue>

	BODY:-
	{
		"cart": {
			"Belts": 3,
			"Trousers": 2,
			"Shoes":1,
			"Shirts":3,
			"Ties":8
		}
	}
	
iii) View Cart

	GET http://localhost:8002/shopping/cart/1234
	
	Output:-
	{
		"ID": 1234,
		"cart": {
			"Belts": 3,
			"Shirts": 3,
			"Shoes": 1,
			"Ties": 8,
			"Trousers": 2
		}
	}
	
iv) Calculate Cart
	
	GET http://localhost:8002/shopping/cart/calculate/1234
	
v)  Update User Cart

	PUT http://localhost:8002/shopping/cart/UpdateUserCart/1234

	Headers:-
	Authorization Bearer <tokenValue>
	
	BODY:-
	{
		"ID": 1234,
		"cart": {
			"Belts": 3,
			"Shirts": 2,
			"Shoes": 1,
			"Ties": 8,
			"Trousers": 2
		}
	}

vi) Delete User Cart
	
	DELETE http://localhost:8002/shopping/cart/DeleteUserCart/1234

	Headers:-
	Authorization Bearer <tokenValue>
	
-------------

Note:- 
i) For basic user addition, no API is available.

ii) As check out is not implemented, concurrency of users is not handled.

iii) For offers, we can come up with a schema to completely automate the whole process in DB updation, for now it is just pipelining

iv) Cart addition need to be in limits of inventory.

v) some corner cases missed










