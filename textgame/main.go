package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ontitansshoulder/textgame/game"
	"github.com/ontitansshoulder/textgame/parser"
	"github.com/pkg/errors"
)

const (
	defaultStoryFile = "gopher.json"
)

func main() {
	err := mainLogic()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func mainLogic() error {
	storyFile := flag.String("story", defaultStoryFile, "the story json file for all story arcs")
	file, err := os.Open(*storyFile)
	if err != nil {
		return errors.Wrap(err, "open story file")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "read story file")
	}

	arcs, err := parser.Parse(bytes)
	if err != nil {
		return err
	}

	newGame, err := game.NewCliGame("intro", arcs)
	if err != nil {
		return errors.Wrap(err, "create cli game")
	}
	newGame.Start()

	return nil
}
