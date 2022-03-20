package user

import (
	"context"
	v "dummystore/lib/variables"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func DeleteUser(user User, authorized bool) (bool, int, error) {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	if authorized{
		err := DeleteUserData(user.ID)
		if err != nil {
			return false, 500, err
		}
		_, err = collection.DeleteOne(ctx, bson.M{"email": user.Email})
		if err != nil {
			return false, 500, v.CouldNotDeleteUser
		}
		return true, 200, nil
	}
	
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return false, 404, v.UserDoesNotExist
	}
	
	valid := CompareHash(user.Password, dbUser.Password)

	if !valid {
		return false, 401, v.UserUnAuthorized
	}

	err = DeleteUserData(dbUser.ID)
	if err != nil {
		return false, 500, err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"email": user.Email})
	if err != nil {
		return false, 500, v.CouldNotDeleteUser
	}

	return true, 200, nil
}


func DeleteUserData(userID primitive.ObjectID) error {
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")
	collection1 := v.Client.Database("Dummystore").Collection("Cart")

	_, err := collection.DeleteMany(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return v.CouldNotDeleteUserProducts
	}

	os.Remove(fmt.Sprintf("./profiles/%s.png", userID.Hex()))

	_, err = collection1.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return	v.CouldNotDeleteUsersCart
	}
	return nil
}