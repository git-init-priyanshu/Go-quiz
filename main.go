package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Problem struct {
	ques string
	ans  string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseProblem(lines [][]string) []Problem {
	problemSlice := make([]Problem, len(lines))
	for i := 0; i < len(lines); i++ {
		problemSlice[i] = Problem{ques: lines[i][0], ans: lines[i][1]}
	}
	return problemSlice
}

func getProblems(filename string) ([]Problem, error) {
	fileObj, err := os.Open(filename)
	if err != nil {
		fmt.Print("Error while opening file")
		return []Problem{}, err
	}

	// Using encoding/csv module to read csv file
	csvReader := csv.NewReader(fileObj)
	linesSlice, err := csvReader.ReadAll() // returns lines in [][]string
	if err != nil {
		fmt.Print("Error while reading lines")
		return []Problem{}, err
	}

	// Parsing lines in defined format
	parsedProblems := parseProblem(linesSlice)
	if err != nil {
		fmt.Print("Erro while parsing lines")
		return []Problem{}, err
	}
	return parsedProblems, nil
}

func main() {
	// Setting up the flags
	fileName := flag.String("f", "questions.csv", "Path of csv file")
	timer := flag.Int("t", 30, "Timer of the quiz")
	flag.Parse()

	// Getting problems from csv file
	problems, err := getProblems(*fileName)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
		log.Fatal(err)
	}

	correctAns := 0

	timeObj := time.NewTimer(time.Duration(*timer) * time.Second)

	ansCh := make(chan string)

problemLooop:
	for i, problem := range problems {
		var answer string
		fmt.Printf("Problem %d: %s = ", i, problem.ques)

		go func() {
			fmt.Scanf("%s", &answer)
			ansCh <- answer
		}()

		select {
		case <-timeObj.C:
			fmt.Print()
			break problemLooop
		case iAns := <-ansCh:
			if iAns == problem.ans {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansCh)
			}
		}
	}

	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
  fmt.Println("Press enter to exit ...")
  <- ansCh
}
