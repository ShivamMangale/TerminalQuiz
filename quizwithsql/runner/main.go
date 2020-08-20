package main

import (
	"TerminalQuiz/quizwithsql/helper"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// var db *sql.DB

func validateProblem(temp problem) string {

	temp.correctoption = strings.ToLower(strings.TrimSpace(temp.correctoption))

	if temp.correctoption != "a" && temp.correctoption != "b" && temp.correctoption != "c" && temp.correctoption != "d" {
		return "Your correctoption isnt a valid letter! It should be either a,b,c or d."
	}

	return "nil"
}

//RunQuiz   The main runner function that takes in the problemset and conducts the quiz.
func RunQuiz(problemSet []problem, timelimit int) (int, []int) {
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
		fmt.Println(i+1, ":  ", p.question, ": \na)", p.optionA, " b)", p.optionB, " c)", p.optionC, " d)", p.optionD)
		fmt.Print("----->")
		go helper.GetAnswer(userAttempt)
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
			if attempt == p.correctoption {
				count++
			} else {
				wrongQuestions = append(wrongQuestions, i+1)
			}
		}
	}

	return count, wrongQuestions
}

//AddQuestion   Given a pointer to a database, one can insert a question with the options.
func AddQuestion(db *sql.DB) {
	fmt.Println("Adding a question!")
	runstat, err := db.Prepare("INSERT INTO problems(question, optionA, optionB, optionC, optionD, correctoption) VALUES (?, ?, ?, ?, ?, ?)")
	helper.CheckError(err)

	fmt.Println("Add the question and give the four options with the correct option when prompted:")
	var temp problem

	temp.question = helper.Inputval("Question")
	temp.optionA = helper.Inputval("Option A")
	temp.optionB = helper.Inputval("Option B")
	temp.optionC = helper.Inputval("Option C")
	temp.optionD = helper.Inputval("Option D")
	temp.correctoption = helper.Inputval("Correct Option(a,b,c or d)")

	if err := validateProblem(temp); err == "nil" {
		fmt.Println("adding : ", temp.question, temp.optionA, temp.optionB, temp.optionC, temp.optionD, temp.correctoption)
		runstat.Exec(temp.question, temp.optionA, temp.optionB, temp.optionC, temp.optionD, temp.correctoption)

		fmt.Println("The question has been added.")
	} else {
		fmt.Println(err)
	}
}

//EditQuestion   Given a pointer to the database, find and edit an already existing question.
func EditQuestion(db *sql.DB) {
	fmt.Println("Editing a question!")
	fmt.Println("Viewing all questions!")
	rows, err := db.Query("SELECT id, question, optionA, optionB, optionC, optionD, correctoption FROM problems")
	helper.CheckError(err)
	var temp problem
	for rows.Next() {
		rows.Scan(&temp.id, &temp.question, &temp.optionA, &temp.optionB, &temp.optionC, &temp.optionD, &temp.correctoption)
		fmt.Println(strconv.Itoa(temp.id) + ": " + temp.question + ":" + temp.optionA + ":" + temp.optionB + ":" + temp.optionC + ":" + temp.optionD + ":" + temp.correctoption) //optionA + optionB
	}

	idtoupdate, err := strconv.Atoi(helper.Inputval("Enter the id of the question you want to edit. (-1 in case you don not want to edit a question)"))
	if err != nil {
		fmt.Println("Enter an integral value please.")
		return
	}
	if idtoupdate == -1 {
		fmt.Println("Chose not to edit. Proceeding.")
		return
	}
	fmt.Println(idtoupdate)

	row := db.QueryRow("SELECT id, question, optionA, optionB, optionC, optionD, correctoption FROM problems WHERE id = (?)", idtoupdate)
	switch err := row.Scan(&temp.id, &temp.question, &temp.optionA, &temp.optionB, &temp.optionC, &temp.optionD, &temp.correctoption); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(temp.id, temp.question, temp.optionA, temp.optionB, temp.optionC, temp.optionD, temp.correctoption)
	default:
		panic(err)
	}

	temp.question = helper.DoUpdate("Question", temp.question)
	temp.optionA = helper.DoUpdate("Option A", temp.optionA)
	temp.optionB = helper.DoUpdate("Option B", temp.optionB)
	temp.optionC = helper.DoUpdate("Option C", temp.optionC)
	temp.optionD = helper.DoUpdate("Option D", temp.optionD)
	temp.correctoption = helper.DoUpdate("Correct Option", temp.correctoption)

	fmt.Println(temp.id, temp.question, temp.optionA, temp.optionB, temp.optionC, temp.optionD, temp.correctoption)

	stmt, err := db.Prepare("UPDATE problems set question=?, optionA=?, optionB=?, optionC=?, optionD=?, correctoption=? where id=?")
	helper.CheckError(err)

	_, err = stmt.Exec(temp.question, temp.optionA, temp.optionB, temp.optionC, temp.optionD, temp.correctoption, idtoupdate)
	helper.CheckError(err)

}

//ViewQuestion   Takes in the pointer to the sql database (*sql.DB) and displays all the questions in the database currently
func ViewQuestion(db *sql.DB) {
	fmt.Println("Viewing all questions!")
	rows, err := db.Query("SELECT id, question, optionA, optionB, optionC, optionD, correctoption FROM problems")
	helper.CheckError(err)
	var temp problem
	for rows.Next() {
		rows.Scan(&temp.id, &temp.question, &temp.optionA, &temp.optionB, &temp.optionC, &temp.optionD, &temp.correctoption)
		fmt.Println(strconv.Itoa(temp.id) + ": " + temp.question + ":" + temp.optionA + ":" + temp.optionB + ":" + temp.optionC + ":" + temp.optionD + ":" + temp.correctoption) //optionA + optionB
	}

}

//StartQuiz   Randomly selects a problemset from the given Quiz database.
func StartQuiz(db *sql.DB) {
	fmt.Println("Starting Quiz")
	rows, err := db.Query("SELECT COUNT(*) FROM problems")
	helper.CheckError(err)
	var val int
	for rows.Next() {
		rows.Scan(&val)
	}

	orderedSet := make([]problem, int(math.Min(float64(val), float64(5))))

	var temp problem

	if val < 5 {
		rows, err := db.Query("SELECT id, question, optionA, optionB, optionC, optionD, correctoption FROM problems")
		helper.CheckError(err)
		ind := 0
		for rows.Next() {
			rows.Scan(&temp.id, &temp.question, &temp.optionA, &temp.optionB, &temp.optionC, &temp.optionD, &temp.correctoption)
			fmt.Println(strconv.Itoa(temp.id) + ": " + temp.question + ":" + temp.optionA + ":" + temp.optionB + ":" + temp.optionC + ":" + temp.optionD + ":" + temp.correctoption) //optionA + optionB
			orderedSet[ind] = temp
			ind++
		}
	} else {
		var done [5]int
		for i := 0; i < 5; i++ {
			done[i] = -1
		}
		var retrieveid int
		for i := 0; i < 5; i++ {
			flag := false
			for flag == false {
				flag = true
				retrieveid = rand.Intn(val-1) + 1
				for _, s := range done {
					if s == retrieveid {
						flag = false
					}
				}
			}
			done[i] = retrieveid
			row := db.QueryRow("SELECT id, question, optionA, optionB, optionC, optionD, correctoption FROM problems WHERE id = (?)", retrieveid)
			switch err := row.Scan(&temp.id, &temp.question, &temp.optionA, &temp.optionB, &temp.optionC, &temp.optionD, &temp.correctoption); err {
			case sql.ErrNoRows:
				fmt.Println("No rows were returned!")
			case nil:
				orderedSet[i] = temp
			default:
				panic(err)
			}
		}
	}
	correct, wrongQuestions := RunQuiz(orderedSet, 60)

	fmt.Println("You got", correct, "of", len(orderedSet), "correct!!!")

	if len(wrongQuestions) != 0 {
		fmt.Print("You got Questions id: ")
		for _, qid := range wrongQuestions {
			fmt.Print(qid, "  ")
		}
		fmt.Println(" wrong. :(")
	}
	fmt.Println("Finito!")
}

func main() {
	fmt.Println("Start!")

	choice := 1
	numoptions := 4
	var dump string
	_ = numoptions

	db, err := sql.Open("sqlite3", "./databank")
	helper.CheckError(err)

	defer db.Close()

	runstat, err := db.Prepare("CREATE TABLE IF NOT EXISTS problems (id INTEGER PRIMARY KEY, question STRING, optionA STRING, optionB STRING, optionC STRING, optionD STRING, correctoption STRING)")
	helper.CheckError(err)

	runstat.Exec()

	for choice != 0 {
		time.Sleep(time.Second / 2)
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		fmt.Println("Menu: ")
		fmt.Println("1: Add a question")
		fmt.Println("2: Edit an already existing question")
		fmt.Println("3: View all questions")
		fmt.Println("4: Start a quiz")
		fmt.Println("0: Quit Session")
		fmt.Println("Enter your choice!(integer choice)")

		fmt.Scanf("%s\n", &dump)
		choice, err = strconv.Atoi(dump)
		if err != nil {
			fmt.Println("Please enter appropriate choice!\n You gave ", dump, " as the choice!")
			choice = numoptions + 10
			continue
		}
		fmt.Println("Your choice was : ", choice)

		switch choice {
		case 0:
			fmt.Println("Thank you for trying it out!")
			break
		case 1:
			AddQuestion(db)
			break
		case 2:
			EditQuestion(db)
			break
		case 3:
			ViewQuestion(db)
			break
		case 4:
			StartQuiz(db)
			break
		default:
			fmt.Println("Please input appropriate choice!")
		}
		helper.GoNext()
	}

	fmt.Println("End!")
}

type problem struct {
	id            int
	question      string
	optionA       string
	optionB       string
	optionC       string
	optionD       string
	correctoption string
}
