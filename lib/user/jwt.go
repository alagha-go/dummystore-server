package user

import (
	"encoding/json"
	"errors"
	v "dummystore/lib/variables"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Token struct {
	Token												string										`json:"token,omitempty"`
	Expires 											time.Time									`json:"expires,omitempty"`
}

var (
	jwtKey = []byte("Password")
)

type Claims struct {
	ID										primitive.ObjectID								`json:"_id,omitempty" bson:"_id,omitempty"`
	UserName								string											`json:"user_name,omitempty" bson:"user_name,omitempty"`
	Email									string											`json:"email,omitempty" bson:"email,omitempty"`
	Seller									bool											`json:"seller,omitempty" bson:"seller,omitempty"`
	RealPassword							string											`json:"real_password,omitempty" bson:"real_password,omitempty"`
	jwt.StandardClaims
}

func GenerateToken(user User) (Token, error) {
	expires := time.Now().Add(72*time.Hour)
	claims  := &Claims{
		ID: user.ID,
		UserName: user.UserName,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		return Token{Token: token, Expires: expires}, errors.New("could not generate authentication token")
	}
	return Token{Token: token, Expires: expires}, nil
}


func ValidateToken(tokenstring string, refresh  ...bool) (bool, interface{}){
	var user User
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenstring, claims, KeyFunc)
	if err != nil {
			return false, nil
	}

	if !tkn.Valid {
		return false, nil
	}else {
		if len(refresh) > 0 {
			if refresh[0] {
				return true, claims
			}
		}
		user.ID = claims.ID
		user.Email = claims.Email
		user.UserName = claims.UserName
		user.NewPassword = claims.RealPassword
		user.Seller = claims.Seller
		if len(user.RealPassword) > 0 {
			user.Public = true
		}

		return true, user
	}
}


func RefreshToken(token string) (Token, error) {
	var user User
	var claims Claims
	valid, data := ValidateToken(token, true)
	if !valid {
		return Token{}, v.InvalidAuthorization
	}
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &claims)

	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 24 * time.Hour {
		return Token{}, v.TokenUnrefreshable
	}
	user.ID = claims.ID
	user.Email = claims.Email
	user.UserName = claims.UserName
	return GenerateToken(user)
}


func KeyFunc(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}