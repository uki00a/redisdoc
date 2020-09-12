package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
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

type Paragraph struct {
	Node
	Nodes []Node
}

type Code struct {
	Node
	Text string
}

type Text struct {
	Node
	Text string
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
	articleMainSelection.Children().Each(func(_ int, s *goquery.Selection) {
		switch {
		case s.Is(".metadata"):
			nodes = append(nodes, *parseMetadata(s))
		case s.Is("p"):
			nodes = append(nodes, *parseParagraph(s))
		default:
			nodes = append(nodes, Text{Text: s.Text()})
		}
	})
	document := &Document{
		CmdName: name,
		Args:    args,
		Article: &Article{nodes},
	}
	return document, nil
}

func parseMetadata(s *goquery.Selection) *Metadata {
	metadata := []string{}
	s.Each(func(_ int, c *goquery.Selection) {
		metadata = append(metadata, strings.TrimSpace(c.Text()))
	})
	return &Metadata{Metadata: metadata}
}

func parseParagraph(s *goquery.Selection) *Paragraph {
	nodes := []Node{}
	s.Children().Each(func(_ int, c *goquery.Selection) {
		switch {
		case c.Is("a"):
			nodes = append(nodes, Text{Text: c.Text()})
		case c.Is("code"):
			nodes = append(nodes, Code{Text: c.Text()})
		default:
			nodes = append(nodes, Text{Text: c.Text()})
		}
	})
	return &Paragraph{Nodes: nodes}
}
