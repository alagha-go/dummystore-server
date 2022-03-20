package dep

import (
	"context"
	v "dummystore/lib/variables"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var (
	Databaseerror = "could not retrieve data from the database"
)

type Department struct {
	ID											primitive.ObjectID					`json:"_id,omitempty" bson:"_id,omitempty"`
	Name										string								`json:"name,omitempty" bson:"name,omitempty"`
	Category									string								`json:"category,omitempty" bson:"category,omitempty"`
	SubCategory									string								`json:"subcategory,omitempty" bson:"subcategory,omitempty"`
	Type										string								`json:"type,omitempty" bson:"type,omitempty"`
	Categories									[]Category							`json:"categories,omitempty" bson:"categories,omitempty"`
}

type Category struct {
	ID											primitive.ObjectID					`json:"_id,omitempty" bson:"_id,omitempty"`
	Name										string								`json:"name,omitempty" bson:"name,omitempty"`
	Url											string								`json:"url,omitempty" bson:"url,omitempty"`
	SubCategories								[]SubCategory						`json:"subcategories,omitempty" bson:"subcategories,omitempty"`
}

type SubCategory struct {
	ID											primitive.ObjectID					`json:"_id,omitempty" bson:"_id,omitempty"`
	Name										string								`json:"name,omitempty" bson:"name,omitempty"`
	Url											string								`json:"url,omitempty" bson:"url,omitempty"`
	Types										[]Type								`json:"types,omitempty" bson:"types,omitempty"`
}

type Type struct {
	ID											primitive.ObjectID					`json:"_id,omitempty" bson:"_id,omitempty"`
	Name										string								`json:"name,omitempty" bson:"name,omitempty"`
	Url											string								`json:"url,omitempty" bson:"url,omitempty"`
}


func GetAllDepartments() ([]Department, int, error){
	var  departments []Department

	ctx := context.Background()
	collection := v.Client.Database("Dummystore").Collection("Departments")
	

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil{
		return nil, 500, v.DatabaseCouldNotRetrieve
	}
		
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var department Department
		cursor.Decode(&department)
		departments = append(departments, department)
	}

	
	return departments, 200, nil
}