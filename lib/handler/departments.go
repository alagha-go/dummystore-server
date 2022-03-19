package handler

import (
	"dummystore/lib/commerce/dep"
	"encoding/json"
	"net/http"
)


func GetAllDepartments(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	data, err, status := dep.GetAllDepartments()
	res.WriteHeader(status)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	encoder := json.NewEncoder(res)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}
