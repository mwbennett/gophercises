package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type question struct {
	prompt string
	answer string
}

func endGame(reason string, score int, totalPossible int) {
	fmt.Println(reason)
	fmt.Printf("Your score: %v/%v\n", score, totalPossible)
	return
}

func confirmReady(timeLimit int) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Are you ready? Press ENTER to start the %v second timer and see the first question.", timeLimit)
	scanner.Scan()
	answer := scanner.Text()

	return answer == ""
}

func loadQuestions(filename *string) ([]*question, error) {
	filePath := fmt.Sprintf("./%v", *filename)
	csvFile, err := os.Open(filePath)
	defer csvFile.Close()

	if err != nil {
		return nil, fmt.Errorf("Unable to open %v, goodbye! %v\n", *filename, err)
	}

	csvReader := csv.NewReader(csvFile)
	rawQuestions, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading questions from csv.")
	}

	var questions []*question

	for _, rawQ := range rawQuestions {
		if len(rawQ) != 2 {
			fmt.Println("Found malformed question")
			continue
		}
		questions = append(questions, &question{
			prompt: rawQ[0],
			answer: rawQ[1],
		})
	}

	return questions, nil
}

func main() {
	// Set up command line flags.
	fileNamePtr := flag.String("csv", "problems.csv", "the file name (with extension) where the questions are located")
	timeLimitPtr := flag.Int("limit", 30, "time limit in seconds")
	flag.Parse()

	questions, err := loadQuestions(fileNamePtr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !confirmReady(*timeLimitPtr) {
		fmt.Println("Exiting, goodbye!")
		return
	}

	totalPossible := len(questions)
	timer := time.NewTimer(time.Duration(*timeLimitPtr) * time.Second)
	score := 0
	endGameReason := "You did it! "

AskQuestions:
	for idx, question := range questions {
		fmt.Printf("Problem #%v: %v\n", idx+1, question.prompt)

		answerChan := make(chan string)
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answerChan <- scanner.Text()
		}()

		select {
		case answer := <-answerChan:
			if answer == question.answer {
				score++
			}
		case <-timer.C:
			endGameReason = "Time's up! "
			break AskQuestions
		}
	}

	endGame(endGameReason, score, totalPossible)
}
