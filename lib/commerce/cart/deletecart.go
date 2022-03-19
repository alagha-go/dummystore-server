package cart

import (
	"context"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func DeleteCart(ID primitive.ObjectID) error {
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}