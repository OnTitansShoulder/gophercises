package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ontitansshoulder/linkparser/analyzer"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

var (
	htmlFileName = flag.String("html", "index.html", "the html file to parse links from")
)

func mainLogic() error {
	flag.Parse()

	htmlFile, err := os.Open(*htmlFileName)
	if err != nil {
		return errors.Wrap(err, "open file")
	}

	doc, err := html.Parse(htmlFile)
	if err != nil {
		return errors.Wrap(err, "parse html")
	}

	linkParser := analyzer.NewLinkParser()
	links, err := linkParser.ParseLink(doc)
	if err != nil {
		return errors.Wrap(err, "parse links")
	}

	linksJson, err := json.MarshalIndent(links, "", "\t")
	if err != nil {
		return errors.Wrap(err, "marshall json")
	}

	fmt.Printf("Parsed Links:\n%s\n", linksJson)
	return nil
}

func main() {
	err := mainLogic()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
