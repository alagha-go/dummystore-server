package handler

import (
	"dummystore/lib/commerce/cart"
	"dummystore/lib/commerce/stats"
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func OrderCart(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}

	id := req.URL.Query().Get("id")
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid id"})
		return
	}

	err = cart.UpdateCart(ID,0, true)
	if err != nil {
		res.WriteHeader(500)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	Cart, status, err := cart.GetCart(ID)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	status, err = stats.AddOrder(Cart.ProductOwnerID, stats.Order{CartID: ID})
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(200)
}


func OrderProduct(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	user, status, err :=  VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	id := req.URL.Query().Get("id")
	quantity, err := strconv.Atoi(req.URL.Query().Get("quantity"))
	if err != nil {
		quantity = 1
	}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid id"})
		return
	}

	_, status, err = cart.AddProductToCart(user.ID, ID, quantity, true)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(200)
}