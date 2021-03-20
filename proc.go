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

func readStatuses(client *mongo.Client, searchtext *string) {
	var status StatusType
	//pnc := make(properNounCounterType, 1000)

	ctx := context.TODO()

	clause1 := primitive.E{Key: "$search", Value: searchtext}
	//clause2 := primitive.E{Key: "$text", Value: clause1}
	// fmt.Printf("clause2 %v", clause2)

	searchfor := bson.D{primitive.E{Key: "$text", Value: bson.D{clause1}}}
	//searchfor := bson.D{}
	findOptions := options.Find()
	findOptions.SetLimit(50000)
	statuses := client.Database("euronews").Collection("statuses")

	cur, err := statuses.Find(ctx, searchfor, findOptions)
	defer cur.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&status); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("status: %v\n", cur.Current)
		fmt.Println("--->", status.Text)
		fmt.Println("***: ", status.Created_at.Time())
	}
}

func estCount(client *mongo.Client) {
	statuses := client.Database("euronews").Collection("statuses")
	opts := options.EstimatedDocumentCount().SetMaxTime(2 * time.Second)
	cnt, err := statuses.EstimatedDocumentCount(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Est. Statuses count: ", cnt)
}

func showProper(cur *mongo.Cursor, dth int) {
	var status StatusType
	pnc := make(properNounCounterType, 1000)

	// ctx := context.TODO()
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&status); err != nil {
			log.Fatal(err)
		}
		matches := pnc.matcher(&status.Text)
		pnc.add(matches)
	}
	pnc.print(dth)
}

func filterStatuses(cur *mongo.Cursor) {
	var status StatusType

	for cur.Next(context.TODO()) {
		if err := cur.Decode(&status); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("status: %v\n", cur.Current)
		fmt.Println("***: ", status.Created_at.Time())
		fmt.Println("--->", status.Text)
	}
}
