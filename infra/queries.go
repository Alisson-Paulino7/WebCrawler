package infra

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VisitedLink struct {
	Website 	string    `bson:"website"`
	Link 		string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_Date"`
}

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
func GetLinks() ([]VisitedLink, error) {

	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	cursor, err := c.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var links []VisitedLink
	for cursor.Next(context.Background()) {
		var link VisitedLink
		cursor.Decode(&link)
		links = append(links, link)
	}

	return links, nil
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