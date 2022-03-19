package products

import (
	"context"
	v "dummystore/lib/variables"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//// this function updates product fields. if a field is empty it will retain the last property
func UpdateOneProduct(product Product) (UpdateResponse, int) {
	var oldProduct Product

	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	err := collection.FindOne(ctx, bson.M{"_id": product.ID}).Decode(&oldProduct)
	if err != nil {
		return UpdateResponse{Product: product, Error: fmt.Errorf("could not find product")}, 404
	}

	if product.ASIN == "" {
		product.ASIN = oldProduct.ASIN
	}

 	if len(product.About) < 1 {
		product.About = oldProduct.About
	}

	if product.Description == "" {
		product.Description = oldProduct.Description
	}

	if product.Price == 0 {
		product.Price = oldProduct.Price
	}

	if product.Currency == "" {
		product.Currency = oldProduct.Currency
	}

	if product.Title == "" {
		product.Title = oldProduct.Title
	}

	if len(product.Sizes) < 1 {
		product.Sizes = oldProduct.Sizes
	}

	product.Department = oldProduct.Department
	product.Rating = oldProduct.Rating

	if len(product.Colors) < 1 {
		product.Colors = oldProduct.Colors
	}

	if product.Brand == "" {
		product.Brand = oldProduct.Brand
	}

	filter := bson.M{"_id": bson.M{"$eq": product.ID}}
	update := bson.M{"$set": product}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UpdateResponse{Product: product, Error: fmt.Errorf("could not update product")}, 500
	}

	return UpdateResponse{Product: product, Success: true}, 200
}
