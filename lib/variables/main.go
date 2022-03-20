package variables

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	Client *mongo.Client
	Router = mux.NewRouter()
	ClientErrors = []string{"invalid Authorization token", "token cannot be refreshed at this time", "user does not exist", "wrong Password", "invalid email address", "token cannot be refreshed at this time", "Unauthorized", "user already exists", "make sure all the fields are filled", "could not find product", "wrong endpoint", }
	ServerErrors = []string{"could not generate authentication token", "could not update user", "could not delete user", "could not update product", }
	UserDoesNotExist = errors.New("user does not exist")
	CartDoesNotExist = errors.New("cart does not exist")
	ProductDoesNotExist = errors.New("product does not exist")
	UserUnAuthorized = errors.New("unauthorized")
	DatabaseCouldNotRetrieve = errors.New("could not retrieve data from the database")
	DatabaseCouldNotSave = errors.New("could not save data to the database")
	InvalidEmail = errors.New("invalid email address")
	UserExists = errors.New("user alrady exists in our atabase")
	CouldNotWriteImageFile = errors.New("could not write you profile image file")
	EmptyFields = errors.New("make sure all fields are filled")
	CouldNotDeleteUser = errors.New("could not delete user")
	CouldNotDeleteUserProducts = errors.New("could not delete user products")
	CouldNotDeleteUsersCart = errors.New("could not delete products in the user's cart")
	InvalidAuthorization = errors.New("invalid authorization token")
	TokenUnrefreshable = errors.New("token cannot be refreshed at this time")
	WrongPassword = errors.New("wrong password")
	CouldNotUpdateUser = errors.New("could not update user")
	CouldNotsaveData = errors.New("could not save data to the database")
	CouldNotUpdateData = errors.New("could not update data")
	NoRelatedProducts = errors.New("could not find products related to this one")
)


func SetErrorHeader(err error, res *http.ResponseWriter) {

}