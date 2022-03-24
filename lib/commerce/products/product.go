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


func AddNewProduct(product Product) (Product, int, error) {
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()
	product.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		return Product{}, 500, v.DatabaseCouldNotSave
	}
	collection.FindOne(ctx, bson.M{"_id": product.ID}).Decode(&product)
	return product, 201, nil
}


func OwnProduct(productID, userID primitive.ObjectID) (Product, int, error) {
	var product Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	err := collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return product, 404, v.ProductDoesNotExist
	}

	if product.OwnerID.Hex() != "000000000000000000000000" {
		return product, 400, errors.New("Product is owned by someone else")
	}

	filter := bson.M{"_id": bson.M{"$eq": product.ID}}
	update := bson.M{"$set": bson.M{"owner_id": userID}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return product, 404, v.CouldNotUpdateData
	}

	collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	return product, 200, nil
}

func GetProductByDepartment(depID primitive.ObjectID) ([]Product, int, error) {
	var products []Product
	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.M{"department": bson.M{"_id": depID}})
	if err != nil {
		return products, 500, v.DatabaseCouldNotRetrieve
	}
	cursor.All(ctx, &products)
	cursor.Close(ctx)

	return products, 200,nil
}