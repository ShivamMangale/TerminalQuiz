package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("The program has started!")

	db, err := sql.Open("sqlite3", "./databank")

	if err != nil {
		panic(err.Error())
	}

	// _ = db

	runstat, err := db.Prepare("CREATE TABLE IF NOT EXISTS QUESTIONS (id INTEGER PRIMARY KEY, question STRING, optionA STRING, optionB STRING)")

	if err != nil {
		panic(err.Error())
	}

	runstat.Exec()

	runstat, err = db.Prepare("INSERT INTO QUESTIONS(question, optionA, optionB) VALUES (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	runstat.Exec("Works?", "Yes", "No")

	rows, err := db.Query("SELECT id, question, optionA, optionB FROM QUESTIONS")

	if err != nil {
		panic(err.Error())
	}

	var id int
	var question string
	var optionA string
	var optionB string

	for rows.Next() {
		rows.Scan(&id, &question, &optionA, &optionB)
		fmt.Println(strconv.Itoa(id) + ": " + question) //optionA + optionB
	}

	fmt.Println("The program has finished!")
}
