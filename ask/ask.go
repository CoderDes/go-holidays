package ask

import (
	"fmt"
	"os"
)

func DefinePart(partNum int, partName string) bool {
	var res bool
	var answer string

	if partNum == 1 {
		fmt.Println("Greetings.")
	}

	fmt.Printf("Let's dive into the part %v of the task -- %v. Would you like to continue? y or n: ", partNum, partName)
	fmt.Fscan(os.Stdin, &answer)

	if isValid := answerValidation(answer); !isValid {
		fmt.Println("Sorry, you entered invalid answer. Please, try again.")
		return DefinePart(partNum, partName)
	}
	if answer == "y" {
		fmt.Println("Nice. Let's go!")
		res = true
	}
	if answer == "n" {
		fmt.Println("Maybe, another time... bye!")
		res = false
	}

	return res
}

func AskCredsToConnect() (string, string, string) {
	var user, pass, db string
	fmt.Print("Enter database user: ")
	fmt.Fscan(os.Stdin, &user)
	fmt.Print("Enter password: ")
	fmt.Fscan(os.Stdin, &pass)
	fmt.Print("Enter database name you want connect to: ")
	fmt.Fscan(os.Stdin, &db)

	return user, pass, db
}

func answerValidation(answ string) bool {
	return answ == "y" || answ == "n"
}
