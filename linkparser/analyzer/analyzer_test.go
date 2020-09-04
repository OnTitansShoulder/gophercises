package analyzer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

const testHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/dog">
    <span>Something in a span</span>
    Text not in a span
    <b>Bold text!</b>
    <a href="Z">Another one!</a>
  </a>
</body>
</html>
`

func TestParseLink(t *testing.T) {
	htmlText := bytes.NewBufferString(testHtml)
	doc, err := html.Parse(htmlText)
	if err != nil {
		t.Error(err)
	}
	linkParser := NewLinkParser()
	links, err := linkParser.ParseLink(doc)
	if err != nil {
		t.Error(err)
	}
	expectedLinks := []Link{
		{
			Href: "/dog",
			Text: "Something in a span Text not in a span Bold text!",
		},
		{
			Href: "Z",
			Text: "Another one!",
		},
	}
	if !cmp.Equal(links, expectedLinks) {
		t.Error(fmt.Sprintf("%v is not equal to %v", links, expectedLinks))
	}
}
