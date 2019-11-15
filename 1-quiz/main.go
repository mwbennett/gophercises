package main

import (
	"bufio"
	"encoding/csv"
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
	text := scanner.Text()

	return question[1] == text
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
	csvFile, err := os.Open("./problems.csv")
	defer csvFile.Close()

	if err != nil {
		fmt.Println("Unable to open problems.csv, goodbye!")
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
