package stats

import (
	"context"
	"dummystore/lib/commerce/cart"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetAllMyOrders(ownerID primitive.ObjectID) ([]cart.Cart, error) {
	var carts []cart.Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")

	cursor, err := collection.Find(ctx, bson.M{"product_owner_id": ownerID, "ordered": true})
	if err != nil {
		return carts, v.DatabaseCouldNotRetrieve
	}
	cursor.All(ctx, &carts)
	cursor.Close(ctx)

	collection = v.Client.Database("Dummystore").Collection("Users")

	for cindex, cart := range carts {
		collection.FindOne(ctx, bson.M{"_id":cart.UserID}).Decode(carts[cindex].User)
	}

	return carts, nil
}