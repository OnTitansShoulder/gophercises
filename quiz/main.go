package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	quizStartTemplate      = "\n============\nQuiz Start\n============\n"
	quizCompletionTemplate = "\n============\nQuiz complete\n============\n"
	questionTemplate       = "Question #%d: %s? Your Answer: "
	completionTemplate     = "Total questions: %d\nTotal attempted: %d\nCorrectly answered: %d\nCorrect rate: %.2f%\n"
)

type quizResult struct {
	Total      int
	Attempted  int
	Correct    int
	Percentage float64
}

func main() {
	err := mainLogic()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func mainLogic() error {
	var problems string
	var shuffle bool
	var timeLimit string
	flag.StringVar(&problems, "problems", "problems.csv", "The file to load the problems")
	flag.BoolVar(&shuffle, "shuffle", false, "Whether to shuffle the question bank.")
	flag.StringVar(&timeLimit, "time-limit", "0", "Time within which to finish the quiz")
	flag.Parse()

	// read in the problems
	problemsFile, err := os.Open(problems)
	if err != nil {
		return errors.Wrapf(err, "open %s", problems)
	}
	defer problemsFile.Close()

	csvReader := csv.NewReader(problemsFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return errors.Wrapf(err, "reading %s", problems)
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	// start the quiz
	fmt.Print(quizStartTemplate)
	cmdChan := make(chan string, 1)
	resultChan := make(chan quizResult)
	errChan := make(chan error)
	defer close(cmdChan)
	defer close(resultChan)
	defer close(errChan)

	if timeLimit != "0" {
		duration, err := time.ParseDuration(timeLimit)
		if err != nil {
			return errors.Wrap(err, "parsing time limit duration")
		}
		go startQuiz(records, cmdChan, resultChan, errChan)
		fmt.Printf("You have %s to do it\n", duration.String())
		time.Sleep(duration)
		fmt.Println("\nTime is up!")
		cmdChan <- "stop"
	} else {
		go startQuiz(records, cmdChan, resultChan, errChan)
	}

	select {
	case result := <-resultChan:
		fmt.Print(quizCompletionTemplate)
		fmt.Printf(completionTemplate, result.Total, result.Attempted, result.Correct, result.Percentage)
	case err := <-errChan:
		return err
	}

	return nil
}

func startQuiz(records [][]string, cmdChan <-chan string, resultChan chan<- quizResult, errChan chan<- error) {
	answerReader := bufio.NewReader(os.Stdin)
	answered := 0
	correct := 0
	for i, record := range records {
		question := record[0]
		correctAnswer := record[1]
		fmt.Printf(questionTemplate, i+1, question)
		answer, err := answerReader.ReadString('\n')
		if err != nil {
			errChan <- errors.Wrap(err, "reading answer")
			return
		}
		select {
		case cmd := <-cmdChan:
			if cmd == "stop" {
				resultChan <- quizResult{
					Total:      len(records),
					Attempted:  answered,
					Correct:    correct,
					Percentage: float64(correct) / float64(len(records)) * 100,
				}
				return
			}
		default:
		}
		answered += 1
		answer = strings.Trim(answer, " \n\r")
		if strings.ToLower(answer) == strings.ToLower(correctAnswer) {
			correct += 1
		}
	}
	resultChan <- quizResult{
		Total:      len(records),
		Attempted:  answered,
		Correct:    correct,
		Percentage: float64(correct) / float64(len(records)) * 100,
	}
}
