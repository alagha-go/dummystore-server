package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func GetProfilePicture(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid user id"})
	}
	res.WriteHeader(200)
	http.ServeFile(res, req, fmt.Sprintf("./profiles/%s.png", id))
}