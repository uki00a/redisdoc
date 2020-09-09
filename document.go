package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type CommandDescription struct {
	Name     string
	Args     []string
	Metadata []string
}

func ParseDocument(body io.Reader) (*CommandDescription, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	cmdSelection := doc.Find(".command")
	name := cmdSelection.Find(".name").Text()
	args := []string{}
	cmdSelection.Find(".arg").Each(func(_ int, selection *goquery.Selection) {
		args = append(args, selection.Text())
	})

	articleMainSelection := doc.Find(".article-main")
	metadata := []string{}
	articleMainSelection.Find(".metadata").Each(func(_ int, selection *goquery.Selection) {
		metadata = append(metadata, selection.Text())
	})
	cmdDescription := &CommandDescription{
		Name:     name,
		Args:     args,
		Metadata: metadata,
	}
	return cmdDescription, nil
}
