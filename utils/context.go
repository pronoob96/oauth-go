package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (c *Context) Close() {
	c.Client.Disconnect(context.TODO())
}

func (c *Context) Collection(name string) *mongo.Collection {
	return c.Database.Collection(name)
}

func NewContext() *Context {
	client, database := GetDBConnection()
	context := &Context{
		Client:   client,
		Database: database,
	}
	return context
}
