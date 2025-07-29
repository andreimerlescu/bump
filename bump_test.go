package main

import (
	"reflect"
	"testing"
)

func Test_correctContents(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := correctContents(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("correctContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("correctContents() got = %v, want %v", got, tt.want)
			}
		})
	}
}
