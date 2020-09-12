package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "expected at least one argument")
		os.Exit(1)
	}

	scraper := NewScraper(http.DefaultClient)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url, err := NewURLFromArgs(args)
	if err != nil {
		log.Fatalf("could not parse URL: %v", err)
	}

	doc, err := scraper.Scrape(ctx, url)
	if err != nil {
		log.Fatalf("could not fetch document: %v", err)
	}

	au := aurora.NewAurora(true)
	printer := NewPrinter(os.Stdout, au, doc)
	printer.Print()
}
