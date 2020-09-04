package analyzer

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type LinkParser struct {
	Links []Link
}

func NewLinkParser() *LinkParser {
	links := make([]Link, 0, 10)
	return &LinkParser{
		Links: links,
	}
}

func (p *LinkParser) ParseLink(doc *html.Node) ([]Link, error) {
	p.findLink(doc)
	return p.Links, nil
}

func (p *LinkParser) findLink(n *html.Node) {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		link := analyzeLink(n)
		p.Links = append(p.Links, link)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.findLink(c)
	}
}

func analyzeLink(n *html.Node) Link {
	href := findAttribute(n.Attr, "href")
	text := findText(n)
	text = strings.TrimSpace(text)
	return Link{
		Href: href,
		Text: text,
	}
}

func findAttribute(attrs []html.Attribute, attrName string) string {
	for _, attr := range attrs {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func findText(n *html.Node) string {
	var texts []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.A {
				continue
			} else {
				texts = append(texts, findText(c))
			}
		} else if c.Type == html.TextNode {
			texts = append(texts, strings.TrimSpace(c.Data))
		}
	}
	return strings.Join(texts, " ")
}
