package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

/*

Flags:
 - filename
 - limit (pt 2)
*/

func ask(number int, question []string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Question %v: %v\n", number, question[0])
	scanner.Scan()
	answer := scanner.Text()

	return question[1] == answer
}

func playQuiz(questions [][]string) (int, int) {
	score := 0

	for idx, question := range questions {
		correct := ask(idx+1, question)
		if correct {
			score++
		}
	}

	return score, len(questions)
}

func main() {
	// Set up command line flags.
	fileNamePtr := flag.String("file_name", "problems.csv", "the file name (with extension) where the questions are located")
	flag.Parse()

	filePath := fmt.Sprintf("./%v", *fileNamePtr)
	csvFile, err := os.Open(filePath)
	defer csvFile.Close()

	if err != nil {
		fmt.Printf("Unable to open problems.csv, goodbye! %v\n", err)
		return
	}

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading rows from csv.")
		return
	}

	correct, possible := playQuiz(rows)

	fmt.Printf("Your score: %v/%v\n", correct, possible)
}
