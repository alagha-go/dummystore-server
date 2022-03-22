package handler

import (
	"dummystore/lib/commerce/stats"
	"encoding/json"
	"net/http"
)

func GetMyStats(res http.ResponseWriter, req *http.Request) {
	user, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	stats := stats.GetMyStatistics(user.ID)
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(stats)
}