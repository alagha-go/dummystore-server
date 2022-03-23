package handler

import (
	p "dummystore/lib/commerce/products"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetProductByAttribute(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var status int
	if req.Method != "GET" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	var data p.Product
	var err error
	id := req.URL.Query().Get("id")
	asin := req.URL.Query().Get("asin")

	if id != "" {
		data, status, err = p.GetProductByID(id)
		res.WriteHeader(status)
	}else if asin != ""{
		data, status, err = p.GetProductByAsin(asin)
		res.WriteHeader(status)
	}else {
		json.NewEncoder(res).Encode(Error{Error: "provide a product id or asin"})
		return
	}

	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}


func UpdateOneProduct(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "PUT" && req.Method != "UPDATE" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	var product p.Product
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: "could not unmarshal data"})
		return
	}
	data, status := p.UpdateOneProduct(product)
	res.WriteHeader(status)
	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}


func AddNewProduct(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	user, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	if !user.Seller {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: "user is not a seller"})
	}
	var product p.Product

	err = json.NewDecoder(req.Body).Decode(&product)

	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: "invalid json"})
	}

	product.OwnerID = user.ID

	product, status, err = p.AddNewProduct(product)
	res.WriteHeader(status)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(product)
}

func OwnProduct(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}

	user, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	if !user.Seller {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: "user is not a seller"})
	}

	id := req.URL.Query().Get("id")
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid id"})
		return
	}

	product, status, err := p.OwnProduct(ID, user.ID)
	res.WriteHeader(status)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	json.NewEncoder(res).Encode(product)
}