package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//Helperplease   is a connection checking variable
var Helperplease = "Connecting!"

func GetAnswer(userAttempt chan string) {
	attempt := ""
	fmt.Scanf("%s\n", &attempt)
	attempt = strings.ToLower(attempt)

	userAttempt <- attempt
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func GoNext() {
	var dummy string
	fmt.Println("Press ENTER to continue.")
	fmt.Scanf("%s\n", &dummy)
}

func Inputval(req string) string {
	fmt.Println(req + ": ")
	inputReader := bufio.NewReader(os.Stdin)
	retval, _ := inputReader.ReadString('\n')
	retval = strings.TrimSuffix(retval, "\n")

	return retval
}

func DoUpdate(term string, val string) string {
	cho := strings.ToLower(Inputval("Do you want to update the " + term + "?(Y or N)"))
	if cho == "y" {
		return Inputval("Enter the " + term + ": ")
	}
	return val
}
