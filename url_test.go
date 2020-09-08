package main

import (
	"fmt"
	"testing"
)

func TestNewURLFromArgs(t *testing.T) {
	var tests = []struct {
		in   []string
		want string
	}{
		{[]string{"GET"}, "https://redis.io/commands/get"},
		{[]string{"set"}, "https://redis.io/commands/set"},
		{[]string{"ACL", "LOAD"}, "https://redis.io/commands/acl-load"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			url, err := NewURLFromArgs(tt.in)
			if err != nil {
				t.Fatal(err)
			}

			if url.String() != tt.want {
				t.Errorf("%s expected, but got %s", tt.want, url.String())
			}
		})
	}
}
