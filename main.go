package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JohnstoneDev/pokedexcli/internal/functions"
)


func main() {
	// All available commands
	commands := functions.GetCommands()

	currentConfig := functions.Config {
		Next: "",
		Previous : "",
	}

	// declare a reader & wait for input
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Pokedex > ")

	// infinite for loop(to emulate a REPL) that waits for inputs & commands
	for {
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input")
			return
		}

		input = strings.TrimSpace(input)

		if cmd, ok := commands[input]; ok {
			err := cmd.Callback(&currentConfig)

			if err != nil {
				fmt.Println(err)
				return
			}

		} else {
			fmt.Println(input, "command not recognized")
		}
		fmt.Print("Pokedex > ")
	}
}
