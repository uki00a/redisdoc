package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type Document struct {
	CmdName string
	Args    []string
	Article *Article
}

type Article struct {
	Nodes []Node
}

type Node interface{}

type Metadata struct {
	Node
	Metadata []string
}

func ParseDocument(body io.Reader) (*Document, error) {
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
	nodes := []Node{}
	articleMainSelection.Children().Each(func(_ int, selection *goquery.Selection) {
		switch {
		case selection.Is(".metadata"):
			metadata := []string{}
			selection.Each(func(_ int, s *goquery.Selection) {
				metadata = append(metadata, s.Text())
			})
			nodes = append(nodes, Metadata{Metadata: metadata})
		default:
			// TODO
		}
	})
	document := &Document{
		CmdName: name,
		Args:    args,
		Article: &Article{nodes},
	}
	return document, nil
}
