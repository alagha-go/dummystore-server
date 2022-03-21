package stats

import (
	"context"
	v "dummystore/lib/variables"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func DeleteUserStats(userID primitive.ObjectID) error {
	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Statistics")

	_, err := collection.DeleteOne(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return errors.New("could not delete user statistics")
	}
	return nil
}