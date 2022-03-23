package cart

import (
	"context"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func UpdateCart(id primitive.ObjectID, quantity int, ordered bool) error {
	var cart Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")

	collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)

	quantity = quantity+cart.Quantity

	if quantity <= 0 {
		quantity = 1
	}

	if cart.Ordered && !ordered{
		ordered = cart.Ordered
	}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": bson.M{
		"quantity": quantity,
		"ordered": ordered,
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return v.CouldNotUpdateData
	}
	return nil
}


func CartExists(id primitive.ObjectID) bool {
	var cart Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)
	return err == nil
}


func ValidateCartOwner(userId, cartId primitive.ObjectID) bool {
	cart, _, _ := GetCart(cartId)
	return cart.UserID == userId
}