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
			fmt.Fprintf(w, "%s%s\n", au.Gray(15, au.Bold("| ")).BgGray(07), au.Gray(19, metadata).BgGray(07))
		}
	case Paragraph:
		for i, node := range node.Nodes {
			if i > 0 {
				// TODO This is not ideal...
				if text, isTextNode := node.(Text); !isTextNode || !StartsWithPunctuation(text.Text) {
					fmt.Fprintf(w, " ")
				}
			}
			p.printNode(node)
		}
		fmt.Fprintf(w, "\n")
	case Example:
		for _, line := range strings.Split(node.Text, "\n") {
			fmt.Fprintf(w, "%s %s\n", au.Green(au.Bold("|")), strings.TrimSpace(line))
		}
	case List:
		for _, item := range node.Items {
			fmt.Fprintf(w, "%s%s\n", au.Green("• "), item)
		}
	case Heading:
		fmt.Fprintf(w, "%s\n", au.Bold(node.Text))
	case Text:
		fmt.Fprintf(w, "%s", strings.TrimSpace(node.Text))
	case Code:
		fmt.Fprintf(w, "%s", au.BgBlue(node.Text))
	}
}
