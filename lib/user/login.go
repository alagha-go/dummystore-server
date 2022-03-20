package user

import (
	"context"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
)


func Login(user User) (Token, int, error) {
	var dbUser User

	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")

	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return Token{}, 404, v.UserDoesNotExist
	}

	valid := CompareHash(user.Password, dbUser.Password)
	if !valid {
		return Token{}, 401, v.WrongPassword
	}

	token, err := GenerateToken(dbUser)
	if err != nil {
		return token, 500, err
	}

	return token, 200, nil
}