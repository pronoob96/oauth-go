package utils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

//GetDBConnection godocs
func GetDBConnection() (*mongo.Client, *mongo.Database) {
	return client, database
}

//CreateDBConnection godocs
func CreateDBConnection(ConnectionURI string, databaseName string) {
	log.Println("Connecting to the following URI: ", ConnectionURI)
	clientOptions := options.Client().ApplyURI(ConnectionURI)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panicln(err)
		return
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panicln(err)
		return
	}

	database = client.Database(databaseName)

	log.Println("Connected to MongoDB! With Database - " + databaseName)
}
