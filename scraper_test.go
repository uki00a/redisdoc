package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type RoundTripFn func(req *http.Request) *http.Response

func (f RoundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newStubHttpClient(fn RoundTripFn) *http.Client {
	return &http.Client{
		Transport: RoundTripFn(fn),
	}
}

func TestScraper(t *testing.T) {
	httpClient := newStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewBufferString(`
<html>
<body>
  <h1 class="command">
    <span class="name">GET</span>
  </h1>
</body>
</html>
      `)),
			Header: make(http.Header),
		}
	})
	scraper := NewScraper(httpClient)
	ctx := context.TODO()
	url, err := url.Parse("https://redis.io/commands/get")
	if err != nil {
		t.Fatal(err)
	}

	result, err := scraper.Scrape(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Errorf("result should not be nil")
	}
	if result.CmdName != "GET" {
		t.Errorf("GET expected, but got %s", result.CmdName)
	}
}
