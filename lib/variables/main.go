package variables

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	Client *mongo.Client
	Router = mux.NewRouter()
	ClientErrors = []string{"invalid Authorization token", "token cannot be refreshed at this time", "user does not exist", "wrong Password", "invalid email address", "token cannot be refreshed at this time", "Unauthorized", "user already exists", "make sure all the fields are filled", "could not find product", "wrong endpoint", }
	ServerErrors = []string{"could not generate authentication token", "could not update user", "could not delete user", "could not update product", }
)


func SetErrorHeader(err error, res *http.ResponseWriter) {

}