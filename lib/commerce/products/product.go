package products

import (
	"context"
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


func GetProductByID(id string) (Product, int, error) {
	var product Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return product, 400, errors.New("invalid id")
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return product, 404, v.ProductDoesNotExist
	}

	return product, 200, nil
}


func GetProductByAsin(asin string) (Product, int, error) {
	var product Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	err := collection.FindOne(ctx, bson.M{"asin": asin}).Decode(&product)
	if err != nil {
		return product, 404, v.ProductDoesNotExist
	}
	
	return product, 200, nil
}