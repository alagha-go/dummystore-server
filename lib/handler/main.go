package handler

import (
	V "dummystore/lib/variables"
	"net/http"
)


func MainHandler() {
	router := V.Router
	
	router.HandleFunc("/", Hello)
	router.HandleFunc("/api/departments", GetAllDepartments)
	router.HandleFunc("/api/updates", GetUpdates)
	router.HandleFunc("/api/products", GetProducts)
	router.HandleFunc("/api/product", GetProductByAttribute)
	router.HandleFunc("/api/signup", CreateUser)
	router.HandleFunc("/api/login", Login)
	router.HandleFunc("/api/refresh", RefreshToken)
	router.HandleFunc("/api/user", GetUser)
	router.HandleFunc("/api/account/delete", DeleteAccount)
	router.HandleFunc("/api/user/delete", DeleteUser)
	router.HandleFunc("/api/cart/add", AddToCart)
	router.HandleFunc("/api/cart/update", UpdateCart)
	router.HandleFunc("/api/cart", GetMyCart)
	router.HandleFunc("/api/profile", ProfileImage)
	router.HandleFunc("/api/profile/{id}", GetProfilePicture)
	router.HandleFunc("/api/orders", GetOrders)
	router.HandleFunc("/api/statistics", GetMyStats)
	router.HandleFunc("/api/cart/delete", DeleteCart)
	router.HandleFunc("/api/product/update", UpdateOneProduct)
	router.HandleFunc("/api/related-products", GetRelatedProducts)
}

func Hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("<h1>Welcome to DummyStore</h1>"))
}


func GetAuth(req *http.Request) (string, error){
	token := req.Header.Get("Authorization")
	if token != "" {
		return token, nil
	}
	cookie, err := req.Cookie("Authorization")
	if err != nil {
		return token, nil
	}
	return cookie.Value, nil
}