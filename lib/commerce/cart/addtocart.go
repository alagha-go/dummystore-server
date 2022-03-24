package cart

import (
	"context"
	"dummystore/lib/commerce/products"
	v "dummystore/lib/variables"
	"errors"
	"os/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Cart struct {
	ID											primitive.ObjectID								`json:"_id,omitempty" bson:"_id,omitempty"`
	ProductOwnerID								primitive.ObjectID								`json:"product_owner_id,omitempty" bson:"product_owner_id,omitempty"`
	Quantity									int												`json:"quantity,omitempty" bson:"quantity,omitempty"`
	ProductID									primitive.ObjectID								`json:"product_id,omitempty" bson:"product_id,omitempty"`
	UserID										primitive.ObjectID								`json:"user_id,omitempty" bson:"user_id,omitempty"`
	User										*user.User										`json:"user,omitempty" bson:"user,omitempty"`
	Product										*products.Product								`json:"product,omitempty" bson:"product,omitempty"`
	Ordered										bool											`json:"ordered,omitempty" bson:"ordered,omitempty"`
	TimeOrdered									time.Time										`json:"time_ordered,omitempty" bson:"time_ordered,omitempty"`
	Paid										bool											`json:"paid,omitempty" bson:"paid,omitempty"`
	Delivered									bool										   	`json:"delivered,omitempty" bson:"delivered,omitempty"`
}


func AddProductToCart(userID, productID primitive.ObjectID, quantity int, ordered bool) (Cart, int, error) {
	var cartexist Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Cart")
	exist, product := ProductExists(productID)
	if !exist {
		return Cart{}, 404, v.ProductDoesNotExist
	}
	if ordered && product.OwnerID.Hex() == "000000000000000000000000" {
		return Cart{}, 400, errors.New("product has no owner")
	}
	cart := Cart{ID: primitive.NewObjectID(), Quantity: quantity,ProductID: productID, UserID: userID, ProductOwnerID: product.OwnerID, Ordered: ordered}

	if ordered {
		cart.TimeOrdered = time.Now()
	}
	
	err := collection.FindOne(ctx, bson.M{"user_id": cart.UserID, "product_id": cart.ProductID}).Decode(&cartexist)
	if err == nil {
		err := UpdateCart(cartexist.ID, quantity, ordered)
		if err != nil {
			return cart, 500, v.ProductDoesNotExist
		}
		return GetCart(cartexist.ID)
	}
	
	_, err = collection.InsertOne(ctx, cart)
	if err != nil {
		return cart, 500, v.CouldNotsaveData
	}
	return cart, 200, nil
}


func ProductExists(id primitive.ObjectID) (bool, products.Product) {
	var product products.Product
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return err == nil, product
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