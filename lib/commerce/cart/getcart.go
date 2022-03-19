package cart

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)




func GetMyCarts(userID primitive.ObjectID) ([]Cart, error) {
	var carts []Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")
	collection1 := v.Client.Database("Dummystore").Collection("Products")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return carts, errors.New("could not retrieve your cart")
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &carts)
	
	for index, cart := range carts {
		_ = collection1.FindOne(ctx, bson.M{"_id": cart.ProductID}).Decode(&carts[index].Product)
	}
	return carts, nil
}