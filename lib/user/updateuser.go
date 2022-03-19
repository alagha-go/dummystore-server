package user

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func UpdateUser(user User) (User, error, int) {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return User{}, errors.New("user does not exist"), 404
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
			return User{}, errors.New("invalid email address"), 400
		}
	}

	if user.NewPassword != "" {
		valid := CompareHash(user.Password, dbUser.Password)
		if !valid {
			return User{}, errors.New("wrong password"), 401
		}
		user.Password = Hasher([]byte(user.NewPassword))
		user.NewPassword = ""
	}

	user.ID = dbUser.ID

	filter := bson.M{"_id": bson.M{"$eq": user.ID}}
	update := bson.M{"$set": user}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return User{}, errors.New("could not update user"), 500
	}

	collection.FindOne(ctx, bson.M{"_id": user.ID}).Decode(&user)
	return user, nil, 200
}


func CheckIfUserExists(userID primitive.ObjectID) bool {
	var dbUser User
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&dbUser)
	return err == nil
}