package stats

import (
	"context"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistics struct {
	OwnerID										primitive.ObjectID					`json:"_id,omitempty" bson:"_id,omitempty"`
	TotalSales									float64								`json:"total_sales,omitempty" bson:"total_sales,omitempty"`
	TotalProducts								int									`json:"total_products,omitempty" bson:"total_products"`
	TotalOrders									int									`json:"total_orders,omitempty" bson:"total_orders,omitempty"`
	PaidOrders									int									`json:"paid_orders,omitempty" bson:"paid_orders,omitempty"`
	SaleStatistics								SaleStatistics						`json:"sale_statistics,omitempty" bson:"sale_statistics,omitempty"`
}

type SaleStatistics struct {
	Years										[]Year								`json:"years,omitempty" bson:"years,omitempty"`					
}

type Year struct {
	Year										int									`json:"year,omitempty" bson:"year,omitempty"`
	Months										[]Month								`json:"months,omitempty" bson:"months,omitempty"`
}

type Month struct {
	Index										int									`json:"index,omitempty" bson:"index,omitempty"`
	Month										string								`json:"month,omitempty" bson:"month,omitempty"`
	Days										[]Day								`json:"days,omitempty" bson:"days,omitempty"`
}

type Day struct {
	Index										int									`json:"index,omitempty" bson:"index,omitempty"`
	Day											string								`json:"day,omitempty" bson:"day,omitempty"`
	Orders										[]Order								`json:"orders,omitempty" bson:"orders,omitempty"`
}

type Order struct {
	ProductID									primitive.ObjectID					`json:"product_id,omitempty" bson:"product_id,omitempty"`
	Paid										bool								`json:"paid,omitempty" bson:"paid,omitempty"`
	Quantity									int									`json:"quantity,omitempty" bson:"quantity,omitempty"`
	BuyserID									primitive.ObjectID					`json:"buyer_id,omitempty" bson:"buyer_id,omitempty"`
}


func GetMyStatistics(userID primitive.ObjectID) (Statistics){
	var stats Statistics
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Statistics")

	collection.FindOne(ctx,   bson.M{"_id": userID}).Decode(&stats)

	return stats
}