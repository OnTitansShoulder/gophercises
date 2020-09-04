package analyzer_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/ontitansshoulder/linkparser/analyzer"
	"golang.org/x/net/html"
)

var _ = Describe("Link Parser", func() {
	DescribeTable("parses links for different cases",
		func(testFile string, expected []analyzer.Link) {
			htmlFile, err := os.Open(testFile)
			Expect(err).Should(BeNil())
			doc, err := html.Parse(htmlFile)
			Expect(err).Should(BeNil())
			linkParser := analyzer.NewLinkParser()
			links, err := linkParser.ParseLink(doc)
			Expect(err).Should(BeNil())
			Expect(links).Should(BeEquivalentTo(expected))
		},
		Entry("parse a simple link", "../test-fixures/single_anchor.html", []analyzer.Link{
			{
				Href: "/other-page",
				Text: "A link to another page",
			},
		}),
		Entry("parse with nested tags", "../test-fixures/nested_tags.html", []analyzer.Link{
			{
				Href: "https://www.twitter.com/joncalhoun",
				Text: "Check me out on twitter",
			},
			{
				Href: "https://github.com/gophercises",
				Text: "Gophercises is on Github !",
			},
		}),
		Entry("parse without comment", "../test-fixures/nested_comment.html", []analyzer.Link{
			{
				Href: "/dog-cat",
				Text: "dog cat",
			},
		}),
		Entry("parse a complex html", "../test-fixures/complex_html.html", []analyzer.Link{
			{
				Href: "#",
				Text: "Login",
			},
			{
				Href: "/lost",
				Text: "Lost? Need help?",
			},
			{
				Href: "https://twitter.com/marcusolsson",
				Text: "@marcusolsson",
			},
		}),
	)
})
