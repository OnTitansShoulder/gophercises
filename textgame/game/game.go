package game

import "github.com/ontitansshoulder/textgame/story"

type TextGame interface {
	Start() error
	Advance(string) story.Option
}

type HtmlGame struct {
}
