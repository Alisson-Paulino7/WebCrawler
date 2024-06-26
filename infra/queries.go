package infra

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// Insert - Insert a document into the database for specific collection type
func Insert(collection string, data interface{}) error {
	// Implementar
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection(collection)

	_, err := c.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	
	return nil
}

// GetLinks - Get a list of links from the database
func FindAllLinks() (estities []VisitedLink, err error) {

	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	opts := options.Find().SetSort(bson.D{{"visited_Date", -1}})

	cursor, err := c.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &estities)

	return
}

// CheckLink - Check if a link already exists in the database
func CheckLink(link string) bool {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	options := options.Count().SetLimit(1)

	n, err := c.CountDocuments(
		context.TODO(),
		bson.D{{Key: "link", Value: link}},
		options,
	)
	if err != nil {
		return false
	}

	return n > 0
}