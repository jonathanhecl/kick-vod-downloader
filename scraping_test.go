package main

import "testing"

func Test_extractVideoID(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test 1", args{"https://kick.com/video/1234"}, "1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVideoID(tt.args.url); got != tt.want {
				t.Errorf("extractVideoID() = %v, want %v", got, tt.want)
			}
		})
	}
}
