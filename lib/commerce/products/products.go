package products

import (
	"context"
	v "dummystore/lib/variables"
	"errors"
	"math/rand"
	"time"

	"dummystore/lib/commerce/dep"

	"github.com/oklog/ulid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Databaseerror = "could not retrieve data from the database"
)

type Response struct {
	NumberOfProducts								int										`json:"total-number-of-products,omitempty"`
	Products										[]Product								`json:"products,omitempty"`
	Token											string									`json:"next_page-token,omitempty"`
}

type Token struct {
	start											int
	Limit											int
	Token											string
	Endpoint										string
	ID												string
}

type Tokens []Token
var tokens Tokens

type Product struct {
	ID												primitive.ObjectID						`json:"_id,omitempty" bson:"_id,omitempty"`
	Title											string									`json:"title,omitempty" bson:"title,omitempty"`
	OwnerID											primitive.ObjectID						`json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	ASIN											string									`json:"asin,omitempty" bson:"asin,omitempty"`
	Price											float64									`json:"price,omitempty" bson:"price,omitempty"`
	Currency										string									`json:"currency,omitempty" bson:"currency,omitempty"`
	Brand											string									`json:"brand,omitempty" bson:"brand,omitempty"`
	Images											[]Image									`json:"images,omitempty" bson:"images,omitempty"`
	Colors											[]Color									`json:"colors,omitempty" bson:"colors,omitempty"`
	Rating											Rate									`json:"rating,omitempty" bson:"rating,omitempty"`
	Sizes											[]string								`json:"sizes,omitempty" bson:"sizes,omitempty"`
	Department										dep.Department								`json:"department,omitempty" bson:"department,omitempty"`
	About											[]string								`json:"about,omitempty" bson:"about,omitempty"`
	Description										string									`json:"description,omitempty" bson:"description,omitempty"`
}


type Image struct {
	Tiny											string									`json:"tiny,omitempty" bson:"tiny,omitempty"`
	Normal											string									`json:"normal,omitempty" bson:"normal,omitempty"`
}

type Rate struct {
	Rating 											float64									`json:"rating,omitempty" bson:"rating,omitempty"`
	Reviews											int										`json:"reviews,omitempty" bson:"reviews,omitempty"`
	Stars											map[int]int								`json:"stars,omitempty" bson:"stars,omitempty"`
}

type Color struct {
	Color											string									`json:"color,omitempty" bson:"color,omitempty"`
	Images											[]Image									`json:"images,omitempty" bson:"images,omitempty"`
}


//// this function gets products as requested by the client
func GetProducts(limit int, token string) (Response, error, int) {
	start := 0
	var response Response
	var products []Product


	collection := v.Client.Database("Dummystore").Collection("Products")
	ctx := context.Background()


/// get all the products from the database
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil{
		return Response{}, errors.New(Databaseerror), 500
	}
		
	defer cursor.Close(ctx)
	cursor.All(ctx, &products)

//////  check if the client has provided a nextpage token
	if token != "" {
		for index, tok := range tokens {
			if tok.Token == token {
				start = tok.start
				if limit == 0 {
					limit = tok.Limit
				}
				tokens = tokens.Remove(index)
			}
		}
	}

///// making sure there is a limit of 5 products if the client provided 0
	if limit == -1 {
		limit = 5
	}

	if limit == 0 {
		limit = len(products)
	}

//// putting the required number of products into the response body
	if start+limit > len(products) {
		limit = len(products) - start
	}
	for index:=start; index < start+limit; index++ {
		response.Products = append(response.Products, products[index])
	}


	response.NumberOfProducts = len(products)

/// generating the next page's token
	var newtoken Token
	newtoken.Token = GenerateToken()
	newtoken.start = start+limit
	newtoken.Limit = limit
	newtoken.Endpoint = "products"
	go RemoveToken(newtoken)
	response.Token = newtoken.Token
	tokens = append(tokens, newtoken)


	return response, nil, 200
}



///// function to remove the used token from the tokens list
func (tokens Tokens) Remove(index int) ([]Token) {
	return append(tokens[:index], tokens[index+1:]...)
}



////// function to generate a unique page token
func GenerateToken() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}


//// func to remove token if it still exists after 3 hours
func RemoveToken(token Token) {
	time.Sleep(3*time.Hour)
	for index, tok := range tokens{
		if token.Token == tok.Token {
			tokens = tokens.Remove(index)
		}
	}
}