package handler

import (
	"dummystore/lib/commerce/stats"
	"encoding/json"
	"net/http"
)

func GetMyStats(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	user, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	if !user.Seller {
		json.NewEncoder(res).Encode(nil)
		return
	}
	stats := stats.GetMyStatistics(user.ID)
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(stats)
}