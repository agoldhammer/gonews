package main

import (
	"testing"
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
		{"test valid srcdb", "euronews", true},
		{"invalid srcdb", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcdb = tt.mysrcdb
			if got := isValidSrcDb(); got != tt.want {
				t.Errorf("isValidSrcDb() = %v, want %v", got, tt.want)
			}
		})
	}
}
