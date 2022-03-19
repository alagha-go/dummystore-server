package user

import (
	"context"
	v "dummystore/lib/variables"
	"errors"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID										primitive.ObjectID								`json:"_id,omitempty" bson:"_id,omitempty"`
	UserName								string											`json:"user_name,omitempty" bson:"user_name,omitempty"`
	Email									string											`json:"email,omitempty" bson:"email,omitempty"`
	Public									bool											`json:"public,omitempty" bson:"public,omitempty"`
	RealPassword							string											`json:"real_password,omitempty" bson:"real_password,omitempty"`
	Password								string											`json:"password,omitempty" bson:"password,omitempty"`
	NewPassword								string											`json:"new_password,omitempty"`
}


func CreateUser(user User) (Token, error, int) {
	var dbUser User
	user.ID = primitive.NewObjectID()

	collection := v.Client.Database("Dummystore").Collection("Users")
	ctx := context.Background()


	if user.Email == "" || user.UserName == "" || user.Password == "" {
		return Token{}, errors.New("make sure all the fields are filled"), 400
	}
	valid := IsEmailValid(user.Email)

	collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if dbUser.Email == user.Email {
		return Token{}, errors.New("user already exists"), 409
	}

	
	if !valid {
		return Token{}, errors.New("invalid email address"), 400
	}
	
	user.Password = Hasher([]byte(user.Password))
	
	token, err := GenerateToken(user)
	if err != nil {
		return token, err, 500
	}
	
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return Token{}, errors.New("could not create user"), 500
	}

	return token, nil, 201
}


func IsEmailValid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}