package main

import (
	"context"
	"fmt"
	"net/http"

	p "dummystore/lib/commerce/products"
	"dummystore/lib/handler"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	port = ":1981"
	mongodb = "mongodb://127.0.0.1:27017/"
)


func main() {

	clientOptions := options.Client().ApplyURI(mongodb)
	ctx := context.Background()
	v.Client, _  = mongo.Connect(ctx, clientOptions)

	fmt.Println("Starting Server...")

	

	go handler.MainHandler()
	go p.UpdateProducts()

	http.ListenAndServe(port, v.Router)
}