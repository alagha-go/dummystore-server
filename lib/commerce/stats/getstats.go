package stats

import (
	"context"
	"dummystore/lib/commerce/cart"
	"dummystore/lib/commerce/products"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistics struct {
	OwnerID										primitive.ObjectID					`json:"owner_id,omitempty" bson:"owner_id,omitempty"`
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
	Day											int									`json:"day,omitempty" bson:"day,omitempty"`
	Orders										[]Order								`json:"orders,omitempty" bson:"orders,omitempty"`
}

type Order struct {
	CartID										primitive.ObjectID					`json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	Cart										interface{}							`json:"cart,omitempty" bson:"cart,omitempty"`
}


func GetMyStatistics(userID primitive.ObjectID) (Statistics){
	var stats Statistics
	var products []products.Product
	var carts	[]cart.Cart
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Statistics")
	collection1 := v.Client.Database("Dummystore").Collection("Products")
	collection2 := v.Client.Database("Dummystore").Collection("Cart")

	collection.FindOne(ctx,   bson.M{"owner_id": userID}).Decode(&stats)

	if stats.OwnerID.Hex() == "000000000000000000000000"{
		CreateNewStatistics(userID)
		return GetMyStatistics(userID)
	}

	cursor, _ := collection1.Find(ctx, bson.M{"owner_id": userID})
	defer cursor.Close(ctx)
	cursor.All(ctx, &products)
	stats.TotalProducts = len(products)

	cursor, _ = collection2.Find(ctx, bson.M{"product_owner_id": userID, "ordered": true, "paid": true})
	cursor.All(ctx, &carts)
	stats.PaidOrders = len(carts)

	for _, Year := range stats.SaleStatistics.Years {
		for _, Month := range Year.Months {
			for _, Day := range Month.Days {
				stats.TotalOrders+=len(Day.Orders)
			}
		}
	}


	return stats
}