package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// command line flags
// source selection
// --srcdb: string, name of source db
// --count: int32
// --start, -s: DateTime
// --end, -e: DateTime
// -m: months ago
// -d: days ago
// -h: hours ago
// --pct: start at p pct of file (neg means from end)
// --search: search string

// output selection
// -o: string: output to screen, db, stream
// --destdb: name of db to store results
// --dth: display threshold

var srcdb string
var count int64
var pct int
var searchfor string
var dth int

func init() {
	flag.StringVar(&srcdb, "srcdb", "", "source db name, e.g. euronews")
	flag.Int64Var(&count, "count", 0, "count (0 or omit for all")
	flag.IntVar(&pct, "pct", 100, "pct of file to proc (neg from end")
	flag.StringVar(&searchfor, "search", "", "search string")
	flag.IntVar(&dth, "dth", 10, "display threshold")
	flag.Parse()
}

func main() {
	println(srcdb)
	println("count", count)
	println("search for: ", searchfor)
	setup()
}

func setup() {
	ctx := context.Background()
	client := connect()
	defer client.Disconnect(ctx)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	// readAuths(client)
	// phrase search example: stext := "\"Emmanuel Macron\""

	statuses := client.Database("euronews").Collection("statuses")
	process(statuses, &searchfor, &count, &dth)
}
