package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"strings"
)

type printer struct {
	w   io.Writer
	au  aurora.Aurora
	doc *Document
}

func NewPrinter(w io.Writer, au aurora.Aurora, doc *Document) *printer {
	return &printer{w, au, doc}
}

func (p *printer) Print() {
	w := p.w
	au := p.au
	doc := p.doc

	fmt.Fprintf(w, "%s %s\n", au.Bold(doc.CmdName), au.Bold(strings.Join(doc.Args, " ")))

	for _, node := range doc.Article.Nodes {
		p.printNode(node)
	}
}

func (p *printer) printNode(node Node) {
	w := p.w
	au := p.au
	switch node := node.(type) {
	case Metadata:
		for _, metadata := range node.Metadata {
			fmt.Fprintf(w, "%s\n", metadata)
		}
	case Paragraph:
		for i, node := range node.Nodes {
			if i > 0 {
				fmt.Fprintf(w, " ")
			}
			p.printNode(node)
		}
		fmt.Fprintf(w, "\n")
	case Text:
		fmt.Fprintf(w, "%s", node.Text)
	case Code:
		fmt.Fprintf(w, "%s", au.BgBlue(node.Text))
	}
}
