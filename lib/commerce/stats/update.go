package stats

import (
	"context"
	"dummystore/lib/commerce/products"
	v "dummystore/lib/variables"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Months = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}


func AddOrder(userID primitive.ObjectID, order Order) (int, error) {
	var stats Statistics

	Year := time.Now().Year()
	monthIndex := time.Now().Month()
	Day := time.Now().Day()

	
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Statistics")
	
	err := collection.FindOne(ctx,   bson.M{"_id": userID}).Decode(&stats)
	if err != nil {
		status, err := CreateNewStatistics(userID)
		if err != nil {
			return status, err
		}
		err = collection.FindOne(ctx,   bson.M{"_id": userID}).Decode(&stats)
		if err != nil {
			return 500, errors.New("unknown error occured")
		}
	}


	for yindex, year := range stats.SaleStatistics.Years {
		if year.Year != Year && yindex < len(stats.SaleStatistics.Years)-1 {
			continue
		}
		if year.Year != Year && yindex == len(stats.SaleStatistics.Years)-1 {
			stats.SaleStatistics.Years =  append(stats.SaleStatistics.Years, CreateNewYear(order))
			break
		}
		for mindex, month := range year.Months {
			if month.Index != int(monthIndex) && mindex < len(year.Months)-1 {
				continue
			}
			if month.Index != int(monthIndex) && mindex == len(year.Months)-1 {
				stats.SaleStatistics.Years[yindex].Months = append(stats.SaleStatistics.Years[yindex].Months, CreateNewMonth(order))
				break
			}
			for dindex, day := range month.Days {
				if day.Day != Day && dindex < len(month.Days)-1 {
					continue
				}
				if day.Day != Day && dindex == len(month.Days)-1 {
					stats.SaleStatistics.Years[yindex].Months[mindex].Days = append(stats.SaleStatistics.Years[yindex].Months[mindex].Days, CreateNewDay(order))
					break
				}
				stats.SaleStatistics.Years[yindex].Months[mindex].Days[dindex].Orders = append(stats.SaleStatistics.Years[yindex].Months[mindex].Days[dindex].Orders, order)
				break
			}
		}
	}

	stats.TotalOrders++

	filter := bson.M{"_id": bson.M{"$eq": stats.OwnerID}}
	update := bson.M{"$set": stats}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 500, errors.New("could not update statistics")
	}

	return 200, nil
}

func CreateNewStatistics(userID primitive.ObjectID) (int, error) {
	var stats Statistics
	var products []products.Product
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Statistics")
	collection1 := v.Client.Database("Dummystore").Collection("Products")

	cursor, err := collection1.Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return 400, v.UserDoesNotExist
	}
	cursor.All(ctx, &products)
	cursor.Close(ctx)
	stats.TotalProducts = len(products)
	stats.OwnerID = userID

	_, err = collection.InsertOne(ctx, stats)
	if err != nil {
		return 500, v.DatabaseCouldNotSave
	}

	return 200, nil
}

func CreateNewYear(order Order) Year {
	return Year{
		Year: time.Now().Year(),
		Months: []Month{CreateNewMonth(order)},
	}
}

func CreateNewMonth(order Order) Month {
	return Month{
		Index: int(time.Now().Month()),
		Month: Months[time.Now().Month()-1],
		Days: []Day{CreateNewDay(order)},
	}
}

func CreateNewDay(order Order) Day {
	return Day{
		Day: time.Now().Day(),
		Orders: []Order{order},
	}
}