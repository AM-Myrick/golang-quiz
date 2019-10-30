package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	csvFilename, shuffle, timeLimit := parseFlags()
	lines := readCSV(csvFilename)
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	scoreQuiz(problems, shuffle, timer)
}

func readCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", filename))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file."))
	}
	return lines
}

func parseFlags() (string, bool, int) {
	csvFilename := flag.String("csv", "randomQuiz.csv", "a csv file in the format of 'question, answer'")
	shuffle := flag.Bool("shuffle", false, "a flag to shuffle the questions")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	numOfQuestions := flag.Int("questions", 20, "the number of questions generated")
	flag.Parse()

	if *csvFilename == "randomQuiz.csv" {
		quizBuilder(*numOfQuestions)
	}

	return *csvFilename, *shuffle, *timeLimit
}

func quizBuilder(quizLength int) {
	var operators = []string{"+", "-", "*"}
	file, err := os.Create("randomQuiz.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to create the CSV file"))
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < quizLength; i++ {
		firstNumber := rand.Intn(20-1) + 1
		secondNumber := rand.Intn(20-1) + 1
		operator := operators[rand.Intn(len(operators))]
		var answer int
		switch operator {
		case "*":
			answer = firstNumber * secondNumber
		case "+":
			answer = firstNumber + secondNumber
		case "-":
			answer = firstNumber - secondNumber
		}
		questionAnswer := fmt.Sprintf("%d%s%d,%d", firstNumber, operator, secondNumber, answer)
		result := []string{questionAnswer}
		writer.Write(result)
	}
}

func scoreQuiz(problems []problem, shuffle bool, timer *time.Timer) {
	if shuffle == true {
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}
	correct := 0
	for i, problem := range problems {
		color.Question.Printf("Problem #%d: %s = ", i+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if strings.ToLower(answer) == problem.answer {
				color.Success.Println("✓")
				correct++
			} else {
				color.Danger.Println("X")
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i, line := range lines {
		if len(line) == 1 {
			line = strings.Split(line[0], ",")
		}
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(strings.ToLower(line[1])),
		}
	}
	return result
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
