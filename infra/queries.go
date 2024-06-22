package infra

import (
	"context"
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