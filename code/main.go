package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func displayErrors(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func readProblems(rows [][]string) []problemstatement {
	sizeOfSet := len(rows)
	orderedSet := make([]problemstatement, sizeOfSet)
	for i, p := range rows {
		orderedSet[i].q = p[0]
		orderedSet[i].a = strings.ToLower(strings.TrimSpace(p[1]))
		orderedSet[i].b = strings.ToLower(strings.TrimSpace(p[2]))
		orderedSet[i].c = strings.ToLower(strings.TrimSpace(p[3]))
		orderedSet[i].d = strings.ToLower(strings.TrimSpace(p[4]))
		orderedSet[i].corr = strings.ToLower(strings.TrimSpace(p[5]))
	}

	return orderedSet
}

func getAnswer(userAttempt chan string) {
	attempt := ""
	fmt.Scanf("%s\n", &attempt)
	attempt = strings.ToLower(attempt)

	userAttempt <- attempt
}

func runQuiz(problemSet []problemstatement, timelimit int) (int, []int) {
	//this function runs the questionnaire, keeping time
	var count int = 0

	fmt.Print("Press ENTER(newline) to start the Quiz. The timer will start instantly.")
	fmt.Println("You will have", timelimit, " seconds to solve the quiz!")
	dummyentry := ""
	fmt.Scanf("%s", &dummyentry)
	timer := time.NewTimer(time.Duration(timelimit) * time.Second)
	userAttempt := make(chan string)
	fmt.Println("Enter your guesses(the letters corresponding to your choice). \nThe questions are as follows:")

	wrongQuestions := make([]int, 0)

	for i, p := range problemSet {
		fmt.Println(i+1, ":  ", p.q, ": \na)", p.a, " b)", p.b, " c)", p.c, " d)", p.d)
		fmt.Print("----->")
		go getAnswer(userAttempt)
		select {
		//this select works as either the time allocated finishes or the user enters attempt
		//if neither happens, it waits
		case <-timer.C:
			//timer returns(interrupts)
			fmt.Println()
			//all the questions unattempted are considered as wrong
			for j := i; j < len(problemSet); j++ {
				wrongQuestions = append(wrongQuestions, j+1)
			}
			return count, wrongQuestions
		case attempt := <-userAttempt:
			//user enters attempt
			if attempt == p.corr {
				count++
			} else {
				wrongQuestions = append(wrongQuestions, i+1)
			}
		}
	}

	return count, wrongQuestions
}

func main() {
	filename := flag.String("csv", "problemset.csv", "Pass the filename of the csv file that contains the problemset for the quiz. \nThe required format is of {question, answer}.")
	timelimit := flag.Int("limit", 60, "Provide a time limit to your quiz attempt to challenge yourself, or ease the speed.\nThe limit is in seconds(s).")
	flag.Parse()

	fileToRead, err := os.Open(*filename) //This is the problemset file.

	if err != nil {
		displayErrors(fmt.Sprintln("Error!!! The file", *filename, "couldnot be accessed. Please recheck.\n", err))
	}

	filePointer := csv.NewReader(fileToRead)
	lines, err := filePointer.ReadAll()
	if err != nil {
		displayErrors(fmt.Sprintln("Could not read the file", *filename, "given"))
	}

	problemSet := readProblems(lines)

	correct, wrongQuestions := runQuiz(problemSet, *timelimit)

	fmt.Println("You got", correct, "of", len(problemSet), "correct!!!")

	if len(wrongQuestions) != 0 {
		fmt.Print("You got Questions id: ")
		for _, qid := range wrongQuestions {
			fmt.Print(qid, "  ")
		}
		fmt.Println(" wrong. :(")
	}
	fmt.Println("Finito!")
}

//format of each question-answer
type problemstatement struct {
	q    string
	a    string
	b    string
	c    string
	d    string
	corr string
}
