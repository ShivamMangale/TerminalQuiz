package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func checkerror(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func gonext() {
	var dummy string
	fmt.Println("Press ENTER to continue.")
	fmt.Scanf("%s\n", &dummy)
}

func inputval(req string) string {
	fmt.Println(req + ": ")
	inputReader := bufio.NewReader(os.Stdin)
	retval, _ := inputReader.ReadString('\n')
	retval = strings.TrimSuffix(retval, "\n")

	return retval
}

func validateProblem(temp problem) string {

	temp.correctoption = strings.ToLower(strings.TrimSpace(temp.correctoption))

	if temp.correctoption != "a" && temp.correctoption != "b" && temp.correctoption != "c" && temp.correctoption != "d" {
		return "Your correctoption isnt a valid letter!"
	}

	return "nil"
}

func getAnswer(userAttempt chan string) {
	attempt := ""
	fmt.Scanf("%s\n", &attempt)
	attempt = strings.ToLower(attempt)

	userAttempt <- attempt
}
