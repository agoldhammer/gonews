package main

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func Test_dateFlag_Set(t *testing.T) {
	var df dateFlag

	err := (&df).Set("3/10/2021")
	want := "2021-03-10"
	if err != nil || want != df.String() {
		t.Errorf("date != 3/10/2021: err: %v, date %v", err, df.String())
	}
	err = (&df).Set("3/10/21")
	if err == nil {
		t.Errorf("bad date 3/10/21  undetected:")
	}
}

func Test_isValidSrcDb(t *testing.T) {
	tests := []struct {
		name    string
		mysrcdb string
		want    bool
	}{
		{"test valid srcdb eu", "euronews", true},
		{"test valid srcdb us", "usnews", true},
		{"invalid srcdb", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.srcdb = tt.mysrcdb
			if got := isValidSrcDb(flags); got != tt.want {
				t.Errorf("isValidSrcDb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDateQuery(t *testing.T) {
	var df = new(dateFlag)
	df.Set("3/17/2021")
	type args struct {
		field string
		op    string
		val   *dateFlag
	}
	tests := []struct {
		name string
		args args
		want bson.E
	}{
		{name: "date query builder", args: args{"created_at", "$gt", df},
			want: bson.E{Key: "created_at", Value: bson.D{{Key: "$gt", Value: df.date}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDateQuery(tt.args.field, tt.args.op, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
