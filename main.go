package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var NEWSURI string = "mongodb://192.168.0.128:27017"

// command line flags
// source selection
// --srcdb: string, name of source db
// --count: int64
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

type dateFlag struct {
	date  time.Time
	valid bool
}

type flagsType struct {
	srcdb     string
	count     int64
	pct       int
	proper    bool
	searchfor string
	dth       int
	start     dateFlag
	end       dateFlag
}

// time-date stuff
var layoutUS = "1/2/2006"
var layoutISO = "2006-01-02"

var flags = new(flagsType)

func main() {
	flag.StringVar(&flags.srcdb, "srcdb", "", "source db name, e.g. euronews")
	flag.Int64Var(&flags.count, "count", 0, "count (0 or omit for all")
	flag.IntVar(&flags.pct, "pct", 100, "pct of file to proc (neg from end")
	flag.BoolVar(&flags.proper, "proper", false, "output condensed proper names")
	flag.StringVar(&flags.searchfor, "search", "", "search string")
	flag.IntVar(&flags.dth, "dth", 10, "display threshold")
	flag.Var(&flags.start, "start", "start date, US fmt e.g. [3/18/2021]")
	flag.Var(&flags.end, "end", "end date, US format e.g. [3/19/2021]")
	flag.Parse()

	printFlags(flags)

	cq := buildCompoundQuery(flags)
	fmt.Println("comp query: ", cq)
	setup(flags)
}

func printFlags(f *flagsType) {
	fmt.Println("Source db: ", f.srcdb)
	fmt.Println("Count: ", f.count)
	fmt.Println("Pct: ", f.pct)
	fmt.Println("Proper: ", f.proper)
	fmt.Println("search for: ", f.searchfor)
	fmt.Println("Display threshold: ", f.dth)
	fmt.Println("Start date: ", f.start.String())
	fmt.Println("End date: ", f.end.String())
}

func isValidSrcDb(f *flagsType) bool {
	switch f.srcdb {
	case
		"euronews",
		"usnews":
		return true
	}
	return false
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

func setup(f *flagsType) {
	ctx := context.Background()
	client := connect()
	defer client.Disconnect(ctx)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Available dbs %v\n---\n", databases)
	if !isValidSrcDb(f) {
		log.Fatal("srcdb must be euronews or usnews")
	}
	// readAuths(client)
	// phrase search example: stext := "\"Emmanuel Macron\""

	statuses := client.Database(f.srcdb).Collection("statuses")
	// TODO: this alternative should be selected from cli
	// filterProperNames(statuses, &f.searchfor, &f.count, &f.dth)
	if cur, err := statusFinder(statuses, f); err != nil {
		log.Fatal("status finder failed")
	} else {
		filterStatuses(cur)
	}

}

// dateFlag satisfied the flag.Value interface
func (df *dateFlag) Set(s string) error {
	date, err := time.Parse(layoutUS, s)
	if err != nil {
		df.valid = false
		return err
	}
	df.date = date
	df.valid = true
	return nil
}

func (df *dateFlag) String() string {
	if df.valid {
		return df.date.Format(layoutISO)
	}
	return "Date not set"
}

func buildTextSearch(flags *flagsType) bson.D {
	searchtext := flags.searchfor
	// clause := bson.E{Key: "$search", Value: searchtext}
	clause := bson.D{{Key: "$search", Value: searchtext}}
	textPart := bson.D{{Key: "$text", Value: clause}}
	return textPart
}

func buildDateQuery(field, op string, val *dateFlag) bson.E {
	if val.valid {
		return bson.E{Key: field, Value: bson.D{{Key: op, Value: val.date}}}
		//return bson.D{{Key: field, Value: bson.D{{Key: op, Value: val.date}}}}
	}
	return bson.E{}
}

func buildCompoundQuery(flags *flagsType) bson.D {
	base := buildTextSearch(flags)
	query := append(base, buildDateQuery("created_at", "$gt", &flags.start))
	return query
}

// if keystr is empty string, "", apply no filter; if limit is 0, apply no limit;
// otherwise, perform text search  on the given collection, coll, according to mongo rules tomatch keystr, which may include quotes
// to seearch for exact phrase. Example: textFinder(statuses, "Macron", 5000)
func statusFinder(coll *mongo.Collection, f *flagsType) (*mongo.Cursor, error) {

	limit := (*flags).count
	//searchfor := buildTextSearch(flags)
	searchfor := buildCompoundQuery(flags)

	fmt.Println("searching for", searchfor)
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	if cur, err := coll.Find(context.TODO(), searchfor, findOptions); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return cur, err
	}
}
