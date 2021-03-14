package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var NEWSURI string = "mongodb://192.168.0.128:27017"

type StatusType struct {
	ID            primitive.ObjectID `bson:"_id"`
	Tid           int                `bson:"id"`
	Author        string             `bson:"author"`
	Created_at    primitive.DateTime `bson:"created_at"`
	Source        string             `bson:"source"`
	Text          string             `bson:"text"`
	Language_code string             `bson:"language_code"`
}

type AuthorType struct {
	ID            primitive.ObjectID `bson:"_id"`
	Author        string             `bson:"author"`
	Language_code string             `bson:"language_code"`
}

func main() {

	ctx := context.Background()
	client := connect()
	defer client.Disconnect(ctx)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	// readAuths(client)
	readStatuses(client)
	fmt.Println(time.Now())
}

// connect to the database at NEWSURI
func connect() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(NEWSURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func readAuths(client *mongo.Client) {
	var author AuthorType
	ctx := context.TODO()
	auths := client.Database("euronews").Collection("authors")
	cur, err := auths.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&author); err != nil {
			panic(err)
		}
		fmt.Printf("author: %v\n", cur.Current)
		fmt.Println("decoded", author)
	}
}

func readStatuses(client *mongo.Client) {
	var status StatusType
	ctx := context.TODO()

	clause1 := primitive.E{Key: "$search", Value: "Castex"}
	// clause2 := primitive.E{Key: "$text", Value: clause1}
	// fmt.Printf("clause2 %v", clause2)

	searchfor := bson.D{primitive.E{Key: "$text", Value: bson.D{clause1}}}
	findOptions := options.Find()
	findOptions.SetLimit(25)
	statuses := client.Database("euronews").Collection("statuses")
	cur, err := statuses.Find(ctx, searchfor, findOptions)
	defer cur.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&status); err != nil {
			panic(err)
		}
		fmt.Printf("status: %v\n", cur.Current)
		fmt.Println("decoded", status)
		fmt.Println("date: ", status.Created_at.Time())
	}
}
