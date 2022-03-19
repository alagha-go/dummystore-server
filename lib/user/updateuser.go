package user

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func UpdateUser(user User) (User, int, error) {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return User{}, 404, errors.New("user does not exist")
	}

	if user.Public {
		user.RealPassword = user.Password
	}


	if user.UserName == "" {
		user.UserName = dbUser.UserName
	}

	if user.Email == "" {
		user.Email = dbUser.Email
	}else {
		valid := IsEmailValid(user.Email)
		if !valid {
			return User{}, 400, errors.New("invalid email address")
		}
	}

	if user.NewPassword != "" {
		valid := CompareHash(user.Password, dbUser.Password)
		if !valid {
			return User{}, 401, errors.New("wrong password")
		}
		user.Password = Hasher([]byte(user.NewPassword))
		user.NewPassword = ""
	}

	user.ID = dbUser.ID

	filter := bson.M{"_id": bson.M{"$eq": user.ID}}
	update := bson.M{"$set": user}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return User{}, 500, errors.New("could not update user")
	}

	collection.FindOne(ctx, bson.M{"_id": user.ID}).Decode(&user)
	return user, 200, nil
}


func CheckIfUserExists(userID primitive.ObjectID) bool {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&dbUser)
	return err == nil
}