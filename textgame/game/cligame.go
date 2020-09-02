package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ontitansshoulder/textgame/story"
)

type CliGame struct {
	InitArc   string
	Arcs      map[string]story.Arc
	BufReader *bufio.Reader
}

func NewCliGame(initArc string, arcs map[string]story.Arc) (*CliGame, error) {
	if _, ok := arcs[initArc]; !ok {
		return nil, fmt.Errorf("initial story arc %s is not defined", initArc)
	}

	cliGame := CliGame{
		InitArc: initArc,
		Arcs:    arcs,
	}
	return &cliGame, nil
}

func (game *CliGame) Start() error {
	currentArc := game.InitArc
	game.BufReader = bufio.NewReader(os.Stdin)
	for {
		option := game.Advance(currentArc)
		if option == story.EmptyOption {
			fmt.Println("End of story.")
			break
		}
		currentArc = option.ArcName
	}
	return nil
}

func (game *CliGame) Advance(arcName string) story.Option {
	currentArc := game.Arcs[arcName]

	printTitle(currentArc.Title)
	printStories(currentArc.Stories)
	if printOptions(currentArc.Options) {
		selection := game.promptSelection(1, len(currentArc.Options))
		return currentArc.Options[selection-1]
	}

	return story.EmptyOption
}

func printTitle(title string) {
	fmt.Printf("###################\n%s\n###################\n", title)
}

func printStories(stories []string) {
	fmt.Println()
	for _, line := range stories {
		fmt.Printf("%s\n\n", line)
	}
}

func printOptions(options []story.Option) bool {
	if len(options) <= 0 {
		return false
	}

	fmt.Println("Here are your options:")

	for i, option := range options {
		i += 1
		fmt.Printf("%d - %s\n", i, option.Text)
	}
	return true
}

func (game *CliGame) promptSelection(floor, cap int) int {
	fmt.Print("\nWhat would you choose?\t")
	var selection int
	for {
		str, err := game.BufReader.ReadString('\n')
		if err != nil {
			fmt.Print("Can't read selection. Try again:\t")
			continue
		}
		str = strings.Trim(str, "\n\r")
		selection, err = strconv.Atoi(str)
		if err != nil {
			fmt.Print("Must enter an integer. Try again:\t")
			continue
		}
		if selection > cap || selection < floor {
			fmt.Print("Enter only available option number. Try again:\t")
			continue
		}
		break
	}
	fmt.Println()
	return selection
}
