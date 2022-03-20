package products

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)




func GetRelatedProductsToID(id, token string, limit int) (Response, int, error){
	var start int = 0
	var newtoken Token
	var product Product
	var products []Product
	var response Response

	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()


	if id == "" && token == "" {
		return response, 400, errors.New("make sure to provide an id or page token")
	}

/// if token is provided by the client make all the fields are field from the token data
	if len(token) > 1 {
		myToken, err := FindToken(token)
		if err != nil {
			return response, 400, err
		}
		if myToken.Endpoint != "productID" {
			return response, 400, errors.New("wrong endpoint")
		}
		start = myToken.start
		id = myToken.ID
		if limit == -1 {
			limit = myToken.Limit
		}
	}

///// if a limit of negative is provided get back default 5

	if limit == -1 {
		limit = 5
	}


	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil && token == ""{
		return response, 400, errors.New("invalid id")
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return response, 404, v.ProductDoesNotExist
	}

	cursor, err := collection.Find(ctx, bson.M{"department": product.Department})
	if err != nil {
		return response, 404, v.NoRelatedProducts
	}

	defer cursor.Close(ctx)
	cursor.All(ctx, &products)

	/// if the client provided a limit of 0 return all the available products
	if limit == 0 {
		limit = len(products)
	}

	response.NumberOfProducts = len(products)

//// looping over the products to append to the response body
	for index:=start; index<start+limit; index++ {
		if index > len(products) {
			break
		}
		response.Products = append(response.Products, products[index])
	}

///// generating a new token for the nextpage
	newtoken.Token = GenerateToken()
	newtoken.start = start+limit
	newtoken.Limit = limit
	newtoken.ID = id
	newtoken.Endpoint = "productID"
	go RemoveToken(newtoken)
	response.Token = newtoken.Token
	tokens = append(tokens, newtoken)
	if token != "" {
		RemoveTok(token)
	}

	return response, 200, nil
}


func FindToken(token string) (Token, error) {
	var newToken Token
	for _, oneToken := range tokens {
		if oneToken.Token == token {
			newToken = oneToken
			return newToken, nil
		}
	}
	return newToken, v.InvalidAuthorization
}


func RemoveTok(token string) {
	for index, oneToken := range tokens {
		if oneToken.Token == token {
			tokens.Remove(index)
			break
		}
	}
}