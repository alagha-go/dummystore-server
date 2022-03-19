package products

import (
	"context"
	v "dummystore/lib/variables"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)


type Data struct {
	AllProducts						int									`json:"allproducts,omitempty"`
	Inserted						int									`json:"inserted,omitempty"`
	Available						int									`json:"available,omitempty"`
	Errors							int									`json:"errors,omitempty"`
}

var updates Data

func UpdateProducts(){
	var products []Product
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &products)
	updates.AllProducts = len(products)


	for _, product := range products {
		time.Sleep(time.Second)
		var newproduct Product
		if len(product.Title) > 0 {
			updates.Available++
			continue
		}
		url := "http://66.228.40.38:5000/product?asin=" + product.ASIN
		res, err := http.Get(url)
		if err != nil {
			_, err := http.Get("http://66.228.40.38:5000/reboot")
			if err != nil {
				log.Panic(err)
			}
		}
		json.NewDecoder(res.Body).Decode(&newproduct)
		product.Title = newproduct.Title
		product.Brand = newproduct.Brand
		product.Price = newproduct.Price
		product.Currency = newproduct.Currency
		product.Images = newproduct.Images
		product.Colors = newproduct.Colors
		product.Sizes = newproduct.Sizes
		product.Rating = newproduct.Rating
		product.Description = newproduct.Description
		product.About = newproduct.About

		err = UpdateProduct(product)
		if err != nil {
			updates.Errors ++
			continue
		}
		updates.Inserted++
	}

}


func UpdateProduct(product Product) error {
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Products")

	filter := bson.M{"_id": bson.M{"$eq": product.ID}}
	update := bson.M{"$set": bson.M{
		"title": product.Title,
		"price": product.Price,
		"currency":product.Currency,
		"brand": product.Brand,
		"images": product.Images,
		"colors": product.Colors,
		"rating": product.Rating,
		"sizes": product.Sizes,
		"about": product.About,
		"description": product.Description,
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetUpdates() Data {
	return updates
}