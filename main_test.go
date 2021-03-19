package main

import "testing"

func Test_dateFlag_Set(t *testing.T) {
	var df dateFlag
	err := (&df).Set("3/10/2021")
	want := "2021-03-10"
	if err != nil || want != df.String() {
		t.Errorf("err: %v, date %v", err, df.String())
	}
}
