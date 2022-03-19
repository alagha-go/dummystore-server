package cart

import (
	"context"
	"dummystore/lib/commerce/products"
	v "dummystore/lib/variables"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Cart struct {
	ID											primitive.ObjectID								`json:"_id,omitempty" bson:"_id,omitempty"`
	Quantity									int												`json:"quantity,omitempty" bson:"quantity,omitempty"`
	ProductID									primitive.ObjectID								`json:"product_id,omitempty" bson:"product_id,omitempty"`
	UserID										primitive.ObjectID								`json:"user_id,omitempty" bson:"user_id,omitempty"`
	Product										interface{}										`json:"product,omitempty" bson:"product,omitempty"`
	Ordered										bool											`json:"ordered,omitempty" bson:"ordered,omitempty"`
	TimeOrdered									time.Time										`json:"time_ordered,omitempty" bson:"time_ordered,omitempty"`
	Paid										bool											`json:"paid,omitempty" bson:"paid,omitempty"`
}


func AddProductToCart(userID, productID primitive.ObjectID, quantity int) (Cart, int, error){
	var cartexist Cart
	cart := Cart{ID: primitive.NewObjectID(), Quantity: quantity,ProductID: productID, UserID: userID}
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")
	exist := ProductExists(productID)
	if !exist {
		return cart, 404, fmt.Errorf("product with the id: %s does not exist", productID.Hex())
	}

	err := collection.FindOne(ctx, bson.M{"user_id": cart.UserID, "product_id": cart.ProductID}).Decode(&cartexist)
	if err == nil {
		err := UpdateCart(cartexist.ID, quantity)
		if err != nil {
			return cart, 500, fmt.Errorf("product with the id: %s does not exist", productID.Hex())
		}
		return GetCart(cartexist.ID)
	}

	_, err = collection.InsertOne(ctx, cart)
	if err != nil {
		return cart, 500, errors.New("could not save product to the cart")
	}
	return cart, 200, nil
}


func ProductExists(id primitive.ObjectID) bool {
	var product products.Product
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return err == nil
}

func GetCart(id primitive.ObjectID) (Cart, int, error) {
	var cart Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)
	if err != nil {
		return cart, 500, err
	}
	return cart, 200, err
}