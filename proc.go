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

func filterProperNames(statuses *mongo.Collection, searchtext *string, count *int64, dth *int) {
	ctx := context.Background()
	if cur, err := textFinder(statuses, searchtext, count); err != nil {
		log.Fatal(err)
	} else {
		showCursor(cur, dth)
		cur.Close(ctx)
	}
	//readStatuses(client, &stext)
	//estCount(client)
	fmt.Println(time.Now())
}

// if keystr is empty string, "", apply no filter; if limit is 0, apply no limit;
// otherwise, perform text search  on the given collection, coll, according to mongo rules tomatch keystr, which may include quotes
// to seearch for exact phrase. Example: textFinder(statuses, "Macron", 5000)
func textFinder(coll *mongo.Collection, keystr *string, limit *int64) (*mongo.Cursor, error) {
	var searchfor bson.D
	// filter := bson.D{{Key: "$bucket", Value: 100}}
	if *keystr != "" {
		clause1 := primitive.E{Key: "$search", Value: *keystr}
		searchfor = bson.D{primitive.E{Key: "$text", Value: bson.D{clause1}}}
	} else {
		searchfor = bson.D{}
	}

	findOptions := options.Find()
	if *limit > 0 {
		findOptions.SetLimit(*limit)
	}

	if cur, err := coll.Find(context.TODO(), searchfor, findOptions); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return cur, err
	}
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
		//matches := pnc.matcher(&status.Text)
		// fmt.Println("matches", matches)
		//pnc.add(matches)
	}
	//pnc.print(100)
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

func showCursor(cur *mongo.Cursor, dth *int) {
	var status StatusType
	pnc := make(properNounCounterType, 1000)

	// ctx := context.TODO()
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&status); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("status: %v\n", cur.Current)
		// fmt.Println("--->", status.Text)
		//fmt.Println("***: ", status.Created_at.Time())
		matches := pnc.matcher(&status.Text)
		// fmt.Println("matches", matches)
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
