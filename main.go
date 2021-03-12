package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var NEWSURI string = "mongodb://192.168.0.128:27017"

func main() {

	ctx, cancel, client := connect()
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	readAuths(client)
	defer client.Disconnect(ctx)
	defer cancel()
}

// connect to the database at NEWSURI
func connect() (ctx context.Context, cancel context.CancelFunc, client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(NEWSURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return ctx, cancel, client
}

func readAuths(client *mongo.Client) {
	auths := client.Database("euronews").Collection("authors")
	cur, err := auths.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.Background()) {
		fmt.Printf("auths: %v\n", cur.Current)
	}
}
