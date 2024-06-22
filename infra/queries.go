package infra

import (
	"context"
	"time"
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

func GetLinks() ([]string, error) {

	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")

	cursor, err := c.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var links []string
	for cursor.Next(context.Background()) {
		var link VisitedLink
		cursor.Decode(&link)
		links = append(links, link.Link)
	}

	return links, nil
}