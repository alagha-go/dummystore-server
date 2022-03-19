package products

import (
	"context"
	"fmt"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type UpdateResponse struct {
	Product										Product									`json:"product,omitempty"`
	Success										bool									`json:"success,omitempty"`
	Error										error									`json:"error,omitempty"`
}


func GetProductByID(id string) (Product, error, int) {
	var product Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return product, errors.New("invalid id"), 400
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return product, fmt.Errorf("could not find a product with the id %s", id), 404
	}

	return product, nil, 200
}


func GetProductByAsin(asin string) (Product, error, int) {
	var product Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	err := collection.FindOne(ctx, bson.M{"asin": asin}).Decode(&product)
	if err != nil {
		return product,fmt.Errorf("could not find a product with the id %s", asin), 404
	}
	
	return product, nil, 200
}