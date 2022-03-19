package handler

import (
	p"dummystore/lib/commerce/products"
	"encoding/json"
	"net/http"
	"strconv"
)


type Error struct {
	Error 										string									`json:"error,omitempty"`
}


func GetProducts(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	var token string
	var limit int
	token = req.URL.Query().Get("token")
	limitstring := req.URL.Query().Get("limit")
	if limitstring != "" {
		limit, _ = strconv.Atoi(limitstring)
	}else {
		limit = -1
	}
	data, status, err := p.GetProducts(limit, token)
	res.WriteHeader(status)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	res.WriteHeader(http.StatusInternalServerError)
	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}

func GetUpdates(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	data := p.GetUpdates()
	json.NewEncoder(res).Encode(data)
}


func GetRelatedProducts(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	var token string
	var id string
	var limit int
	token = req.URL.Query().Get("token")
	id = req.URL.Query().Get("id")
	limitstring := req.URL.Query().Get("limit")
	if limitstring != "" {
		limit, _ = strconv.Atoi(limitstring)
	}else {
		limit = -1
	}

	data, status, err := p.GetRelatedProductsToID(id, token, limit)
	res.WriteHeader(status)
	if err != nil {
		encoder := json.NewEncoder(res)
		encoder.SetEscapeHTML(false)
		encoder.Encode(Error{Error: err.Error()})
		return
	}

	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}