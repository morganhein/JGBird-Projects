package main

// TODO - convert the loaded file into a printable list

/*
 @jg: can you turn the actual functionality into a pkg, and have main() just call/use the package?
 I personally prefer main entrypoints and the namespace to be as slim as possible. For example, the main
 for a microservice i'm writing that is around 12k LOC just has:
func main() {
	l := config.NewLogger(config.LoggerConfig{})

	if err := app.Run(l); err != nil {
		l.Fatal().Msgf("error running: %v", err)
	}
}

that's the entire main function! and it connects to a db, starts up a gprc server, configures authentication, etc.
Keep main small!
*/

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// @jg: comments should start with the function name, so this should be `ReadFile ...`

// This will open a file, convert the text to a slice of strings,
// return the slice, then close the file.
func ReadFile(fString string) []string {
	rFile, err := os.Open(fString)
	if err != nil {
		// @jg: in general, don't panic like this, except in main. Bubble the error up, and let the caller handle it
		log.Fatal(err)
	}
	defer rFile.Close()
	scanner := bufio.NewScanner(rFile)
	var rLines []string
	for scanner.Scan() {
		rLines = append(rLines, scanner.Text())
	}

	return rLines
}

// This will open a file, save the contents as a string,
// And sort the string into a legible list.
func LoadSave(lString string) {
	lFile, err := os.Open(lString)
	if err != nil {
		log.Fatal(err)
	}
	defer lFile.Close()
	scanner := bufio.NewScanner(lFile)
	var lLines []string
	for scanner.Scan() {
		lLines = append(lLines, scanner.Text())
	}
	fmt.Printf("You have %d saved characters:\n", len(lLines))
	for _, line := range lLines {
		lineSlice := strings.Split(line, ",")
		fmt.Printf("%s, a(n) %s %s\n", lineSlice[0], lineSlice[1], lineSlice[2])
	}

}

// This will create a save file to store the character's name,
// race, and class.
func SaveFile(name, race, class string) error {
	// @jg: this isn't os-agnostic, and only runs on windows. If you write an application that uses a different path separator (like unix) this won't
	// work. See if you can figure out how to make a path that is os-agnostic. The stdlib should help with that
	// @jg: why is this called err1? it's idiomatic, unless you need to keep references to multiple errors, to just name it err and
	// override later usages of err.
	sFile, err1 := os.OpenFile("C:\\projects\\personal\\Character Creator\\saves.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err1 != nil {
		log.Fatal(err1)
		// @jg: this return err is never reached, can you guess why?
		return err1
	}

	defer sFile.Close()

	// @jg: what is the better way to concatenate strings? this is not idiomatic go
	// @jg: err2? maybe just re-use err?
	_, err2 := sFile.WriteString(name + "," + race + "," + class + "\n")
	if err2 != nil {
		log.Fatal(err2)
	}
	return err2
}

// This selects a random entry from a given slice of strings
func RandomString(iString []string) string {
	// @jg: this function could be reduced by 1 line, can you see how?
	oString := iString[rand.Intn(len(iString))]
	return oString
}

// @jg: this function is big, it needs to be broken up into smaller functions
func main() {
	rand.Seed(time.Now().UnixNano())
	var charClass, charRace string
	var classList, raceList []string
	var saveErr error
	var menuInput, nameInput string
	inputScanner := bufio.NewScanner(os.Stdin)
	// @jg: this nested loop is too complex. If you break them up into functions, then you don't need the loop identifiers
	// and it makes it easier to read. Try doing that.
MenuLoop:
	for {
		fmt.Println("\n\nCharacter Creation Options")
		fmt.Println("--------------------------")
		fmt.Println("1) Dungeons & Dragons 5th Edition")
		fmt.Println("2) Show saved characters")
		fmt.Println("3) Quit")
		fmt.Println("")
		fmt.Print("Please enter your choice: ")
		inputScanner.Scan()
		menuInput = inputScanner.Text()
		switch menuInput {
		case "1":
			// @jg: os specific pathing....
			classList = ReadFile("C:\\projects\\personal\\Character Creator\\classes.txt")
			raceList = ReadFile("C:\\projects\\personal\\Character Creator\\races.txt")

			charClass = RandomString(classList)
			charRace = RandomString(raceList)
		case "2":
			LoadSave("C:\\projects\\personal\\Character Creator\\saves.txt")
			continue MenuLoop
		case "3":
			break MenuLoop
		default:
			fmt.Print("Invalid entry. ")
			continue MenuLoop
		}

		//displays the character traits and allows the user to name them
	NameLoop:
		for {
			fmt.Printf("\nYour random character is a(n) %v %v. Please name your character: ", charRace, charClass)
			inputScanner.Scan()
			nameInput = inputScanner.Text()
			blankName := strings.TrimSpace(nameInput) == ""
			if blankName {
				fmt.Println("Please enter a name.")
				continue NameLoop
			}

			fmt.Print("Saving...")
			saveErr = SaveFile(nameInput, charRace, charClass)

			if saveErr != nil {
				fmt.Println("Something went wrong. Details below. Please try again.")
				fmt.Println(saveErr)
				continue NameLoop
			} else {
				fmt.Println("Success!")
				break NameLoop
			}
		}
	}
	fmt.Println("Shutting down. Have a nice day!")
}
