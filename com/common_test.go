package com

import (
	"os"
	"testing"
)

func TestPropS(t *testing.T) {
	os.Setenv("KEY_1", "VAL_1")
	type args struct {
		key string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Should override", args{"KEY_1", "NOPE"}, "VAL_1"},
		{"Should take default", args{"KEY_2", "VAL_1"}, "VAL_1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PropS(tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("PropS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJsonString(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Should Serialize", args{map[string]string{"a": "b"}}, `{"a":"b"}`},
		{"Should Return Blank", args{nil}, `null`},
		{"Should Return Blank", args{``}, ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJsonString(tt.args.obj); got != tt.want {
				t.Errorf("ToJsonString() = %v, want %v", got, tt.want)
			}
		})
	}
}
