package products

import (
	"context"
	v "dummystore/lib/variables"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)




func GetRelatedProductsToID(id, token string, limit int) (Response, error, int){
	var start int = 0
	var newtoken Token
	var product Product
	var products []Product
	var response Response

	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()


	if id == "" && token == "" {
		return response, errors.New("make sure to provide an id or page token"), 400
	}

/// if token is provided by the client make all the fields are field from the token data
	if len(token) > 1 {
		myToken, err := FindToken(token)
		if err != nil {
			return response, err, 400
		}
		if myToken.Endpoint != "productID" {
			return response, errors.New("wrong endpoint"), 400
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
		return response, errors.New("invalid id"), 400
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return response, fmt.Errorf("could not find a product with the id: %s", id), 404
	}

	cursor, err := collection.Find(ctx, bson.M{"department": product.Department})
	if err != nil {
		return response, errors.New("could not find products related to this product"), 404
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

	return response, nil, 200
}


func FindToken(token string) (Token, error) {
	var newToken Token
	for _, oneToken := range tokens {
		if oneToken.Token == token {
			newToken = oneToken
			return newToken, nil
		}
	}
	return newToken, errors.New("invalid token")
}


func RemoveTok(token string) {
	for index, oneToken := range tokens {
		if oneToken.Token == token {
			tokens.Remove(index)
			break
		}
	}
}