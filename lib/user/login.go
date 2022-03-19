package user

import (
	"context"
	v "dummystore/lib/variables"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)


func Login(user User) (Token, error, int) {
	var dbUser User

	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Users")

	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return Token{}, fmt.Errorf("user does not exist"), 404
	}

	valid := CompareHash(user.Password, dbUser.Password)
	if !valid {
		return Token{}, errors.New("wrong Password"), 401
	}

	token, err := GenerateToken(dbUser)
	if err != nil {
		return token, err, 500
	}

	return token, nil, 200
}