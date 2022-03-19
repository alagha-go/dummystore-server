package handler

import (
	"dummystore/lib/commerce/cart"
	u "dummystore/lib/user"
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func AddToCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(405)
		return
	}
	var user u.User
	quantity, err := strconv.Atoi(req.URL.Query().Get("quantity"))
	if err != nil {
		quantity = 1
	}
	if quantity == 0 {
		quantity = 1
	}
	user, err, status := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	id := req.URL.Query().Get("id")
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid id"})
		return
	}

	cart, err, status := cart.AddProductToCart(user.ID, ID, quantity)
	res.WriteHeader(status)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	json.NewEncoder(res).Encode(cart)
}

func UpdateCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "PUT" && req.Method != "UPDATE" {
		res.WriteHeader(405)
		return
	}
	
	user, err, status := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	
	id := req.URL.Query().Get("id")
	quantity, err := strconv.Atoi(req.URL.Query().Get("quantity"))
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "provide a valid quantity"})
		return
	}
	
	ID, err := primitive.ObjectIDFromHex(id)
	
	exist := cart.CartExists(ID)
	if !exist {
		res.WriteHeader(404)
		json.NewEncoder(res).Encode(Error{Error: "cart does not exist"})
		return
	}
	valid := cart.ValidateCartOwner(user.ID, ID)

	if !valid {
		res.WriteHeader(401)
		json.NewEncoder(res).Encode(Error{Error: "this cart does not belong to this user"})
		return
	}

	if err != nil {
		 res.WriteHeader(400)
		 json.NewEncoder(res).Encode(Error{Error: "provide a valid id"})
		 return
	}

	err = cart.UpdateCart(ID, quantity)
	if err != nil {
		res.WriteHeader(500)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	cart, err, status := cart.GetCart(ID)
	res.WriteHeader(status)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	json.NewEncoder(res).Encode(cart)
}


func DeleteCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "DELETE" {
		res.WriteHeader(405)
		return
	}

	user, err, status := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	id := req.URL.Query().Get("id")

	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		res.WriteHeader(401)
		json.NewEncoder(res).Encode(Error{Error: "invalid cart id"})
	}

	exist := cart.CartExists(ID)
	if !exist {
		res.WriteHeader(404)
		json.NewEncoder(res).Encode(Error{Error: "cart does not exist"})
		return
	}
	valid := cart.ValidateCartOwner(user.ID, ID)

	if !valid {
		res.WriteHeader(401)
		json.NewEncoder(res).Encode(Error{Error: "this cart does not belong you"})
		return
	}

	err = cart.DeleteCart(ID)
	if err != nil {
		res.WriteHeader(500)
		json.NewEncoder(res).Encode(Error{Error: "could not delete the cart"})
		return
	}
	res.WriteHeader(200)
}


func GetMyCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	user, err, status := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	carts, err := cart.GetMyCarts(user.ID)
	if err != nil {
		res.WriteHeader(500)
		json.NewEncoder(res).Encode(Error{Error: "could not get your cart"})
		return
	}

	res.WriteHeader(200)
	json.NewEncoder(res).Encode(carts)
}