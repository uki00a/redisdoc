package main

import (
	"fmt"
	"io"
	"strings"
)

type printer struct {
	w   io.Writer
	doc *Document
}

func NewPrinter(w io.Writer, doc *Document) *printer {
	return &printer{w, doc}
}

func (p *printer) Print() {
	w := p.w
	doc := p.doc

	fmt.Fprintf(w, "%s %s\n", doc.CmdName, strings.Join(doc.Args, " "))

	for _, node := range doc.Article.Nodes {
		switch node := node.(type) {
		case Metadata:
			for _, metadata := range node.Metadata {
				fmt.Fprintf(w, "%s\n", metadata)
			}
		}
	}
}
