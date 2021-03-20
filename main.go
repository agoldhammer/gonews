package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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

var srcdb string
var count int64
var pct int
var searchfor string
var dth int

type dateFlag struct {
	date  time.Time
	valid bool
}

var start dateFlag
var end dateFlag

// time-date stuff
var layoutUS = "1/2/2006"
var layoutISO = "2006-01-02"

func init() {
	// kludge: see https://github.com/golang/go/issues/31859
	testing.Init()
	flag.StringVar(&srcdb, "srcdb", "", "source db name, e.g. euronews")
	flag.Int64Var(&count, "count", 0, "count (0 or omit for all")
	flag.IntVar(&pct, "pct", 100, "pct of file to proc (neg from end")
	flag.StringVar(&searchfor, "search", "", "search string")
	flag.IntVar(&dth, "dth", 10, "display threshold")
	flag.Var(&start, "start", "start date, US fmt e.g. [3/18/2021]")
	flag.Var(&end, "end", "end date, US format e.g. [3/19/2021]")
	flag.Parse()
}

func main() {
	fmt.Println("Source db: ", srcdb)
	fmt.Println("Count: ", count)
	fmt.Println("Pct: ", pct)
	fmt.Println("search for: ", searchfor)
	fmt.Println("Display threshold: ", dth)
	fmt.Println("Start date: ", start.String())
	fmt.Println("End date: ", end.String())
	//Test_Set()
	setup()
}

func isValidSrcDb() bool {
	switch srcdb {
	case
		"euronews",
		"usnews":
		return true
	}
	return false
}

func setup() {
	ctx := context.Background()
	client := connect()
	defer client.Disconnect(ctx)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Available dbs %v\n---\n", databases)
	if !isValidSrcDb() {
		log.Fatal("srcdb must be euronews or usnews")
	}
	// readAuths(client)
	// phrase search example: stext := "\"Emmanuel Macron\""

	statuses := client.Database(srcdb).Collection("statuses")
	process(statuses, &searchfor, &count, &dth)
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
