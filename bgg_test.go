package main

import (
	"reflect"
	"testing"
)

func Test_getSearchXML(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want Items
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSearchXML(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSearchXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
