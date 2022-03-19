package handler

import (
	p"dummystore/lib/commerce/products"
	"encoding/json"
	"net/http"
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
	if req.Method != "PUT" || req.Method != "UPDATE" {
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