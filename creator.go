package main

// TODO - convert the loaded file into a printable list

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"strings"
)

// This will open a file, convert the text to a slice of strings,
// return the slice, then close the file.
func ReadFile(fString string) []string {
	rFile, err := os.Open(fString)
	if err != nil {
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
func LoadSave(lString string){
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
	for _, line := range lLines{
		lineSlice := strings.Split(line, ",")
		fmt.Printf("%s, a(n) %s %s\n", lineSlice[0], lineSlice[1], lineSlice[2])
	}

}

// This will create a save file to store the character's name,
// race, and class.
func SaveFile(name, race, class string) error {
	sFile, err1 := os.OpenFile("C:\\projects\\personal\\Character Creator\\saves.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err1 != nil {
		log.Fatal(err1)
		return err1
	}

	defer sFile.Close()

	_, err2 := sFile.WriteString(name + "," + race + "," + class + "\n")
	if err2 != nil {
		log.Fatal(err2)
	}
	return err2
}

// This selects a random entry from a given slice of strings
func RandomString(iString []string) string {
	oString := iString[rand.Intn(len(iString))]
	return oString
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var charClass, charRace string
	var classList, raceList []string
	var saveErr error
	var menuInput, nameInput string
	inputScanner := bufio.NewScanner(os.Stdin)

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
