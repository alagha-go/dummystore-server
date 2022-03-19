package user

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func DeleteUser(user User, authorized bool) (bool, error, int) {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	if authorized{
		err := DeleteUserData(user.ID)
		if err != nil {
			return false, err, 500
		}
		_, err = collection.DeleteOne(ctx, bson.M{"email": user.Email})
		if err != nil {
			return false, errors.New("could not delete user"), 500
		}
		return true, nil, 200
	}
	
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return false, errors.New("user does not exist"), 404
	}
	
	valid := CompareHash(user.Password, dbUser.Password)

	if !valid {
		return false, errors.New("Unauthorized"), 401
	}

	err = DeleteUserData(dbUser.ID)
	if err != nil {
		return false, err, 500
	}

	_, err = collection.DeleteOne(ctx, bson.M{"email": user.Email})
	if err != nil {
		return false, errors.New("could not delete user"), 500
	}

	return true, nil, 200
}


func DeleteUserData(userID primitive.ObjectID) error {
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")
	collection1 := v.Client.Database("Dummystore").Collection("Cart")

	_, err := collection.DeleteMany(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return errors.New("could not delete products")
	}

	_, err = collection1.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return	errors.New("could not delet products in your cart")
	}
	return nil
}