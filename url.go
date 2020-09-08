package main

import (
	"fmt"
	"net/url"
	"strings"
)

func NewURLFromArgs(args []string) (*url.URL, error) {
	pageName := strings.ToLower(strings.Join(args, "-"))
	rawURL := fmt.Sprintf("https://redis.io/commands/%s", pageName)
	return url.Parse(rawURL)
}
