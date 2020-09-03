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
	defaultGameMode  = "web"
)

func main() {
	err := mainLogic()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func mainLogic() error {
	storyFile := flag.String("story", defaultStoryFile, "the story json file for all story arcs")
	gameMode := flag.String("mode", defaultGameMode, "the game mode: web/cli")
	flag.Parse()

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
		return errors.Wrap(err, "parse story arcs")
	}

	var newGame game.TextGame
	switch *gameMode {
	case "web":
		newGame, err = game.NewWebGame("intro", arcs)
	case "cli":
		newGame, err = game.NewCliGame("intro", arcs)
	default:
		err = errors.New("unsupported game mode: " + *gameMode)
	}
	if err != nil {
		return errors.Wrap(err, "create game")
	}

	err = newGame.Start()
	if err != nil {
		return err
	}

	return nil
}
