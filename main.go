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
	csvFilename := flag.String("csv", "randomQuiz.csv", "a csv file in the format of 'question, answer'")
	shuffle := flag.Bool("shuffle", false, "a flag to shuffle the questions")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	if *csvFilename == "randomQuiz.csv" {
		quizBuilder(100)
	}

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file."))
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	scoreQuiz(problems, *shuffle, timer)
}

func quizBuilder(numOfQuestions int) {
	var operators = []string{"+", "-", "*"}
	file, err := os.Create("randomQuiz.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to create the CSV file"))
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < numOfQuestions; i++ {
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
				color.Success.Println("Correct!")
				correct++
			} else {
				color.Danger.Println("Keep trying!")
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
