package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

/*

Flags:
 - filename
 - limit (pt 2)
*/

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

func main() {
	// Set up command line flags.
	fileNamePtr := flag.String("file_name", "problems.csv", "the file name (with extension) where the questions are located")
	timeLimitPtr := flag.Int("time_limit", 30, "time limit in seconds")
	flag.Parse()

	filePath := fmt.Sprintf("./%v", *fileNamePtr)
	csvFile, err := os.Open(filePath)
	defer csvFile.Close()

	if err != nil {
		fmt.Printf("Unable to open problems.csv, goodbye! %v\n", err)
		return
	}

	csvReader := csv.NewReader(csvFile)
	questions, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading questions from csv.")
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
			fmt.Printf("Question %v: %v\n", idx + 1, question[0])

			answerChan := make(chan string)
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				answerChan <- scanner.Text()
			}()

			select {
			case answer := <-answerChan:
				if answer == question[1] {
					score++
				}
			case <-timer.C:
				endGameReason = "Time's up! "
				break AskQuestions
			}
		}

	endGame(endGameReason, score, totalPossible)
}
