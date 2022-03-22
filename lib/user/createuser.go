package user

import (
	"context"
	v "dummystore/lib/variables"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID										primitive.ObjectID								`json:"_id,omitempty" bson:"_id,omitempty"`
	UserName								string											`json:"user_name,omitempty" bson:"user_name,omitempty"`
	Email									string											`json:"email,omitempty" bson:"email,omitempty"`
	ImageFile								multipart.File									`json:"image_file,omitempty" bson:"image_file,omitempty"`
	Public									bool											`json:"public,omitempty" bson:"public,omitempty"`
	Seller									bool											`json:"seller,omitempty" bson:"seller,omitempty"`
	RealPassword							string											`json:"real_password,omitempty" bson:"real_password,omitempty"`
	Password								string											`json:"password,omitempty" bson:"password,omitempty"`
	NewPassword								string											`json:"new_password,omitempty" bson:"new_password,omitempty"`
}


func CreateUser(user User) (Token, int, error) {
	var dbUser User
	user.ID = primitive.NewObjectID()

	collection := v.Client.Database("Dummystore").Collection("Users")
	ctx := context.Background()


	if user.Email == "" || user.UserName == "" || user.Password == "" {
		return Token{}, 400, v.EmptyFields
	}
	valid := IsEmailValid(user.Email)

	collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if dbUser.Email == user.Email {
		return Token{}, 409, v.UserExists
	}

	
	if !valid {
		return Token{}, 400, v.InvalidEmail
	}

	if user.Public {
		user.RealPassword = user.Password
	}
	
	user.Password = Hasher([]byte(user.Password))
	
	token, err := GenerateToken(user)
	if err != nil {
		return token, 500, err
	}

	out, _ := os.Create(fmt.Sprintf("./profiles/%s.png", user.ID.Hex()))

	_, err = io.Copy(out, user.ImageFile)

	if err != nil {
		return token, 500, v.CouldNotWriteImageFile
	}

	user.ImageFile = nil
	
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return Token{}, 500, v.DatabaseCouldNotSave
	}

	return token, 201, nil
}


func IsEmailValid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}