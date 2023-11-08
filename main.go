package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/JohnstoneDev/pokedexcli/internal/functions"
	"github.com/JohnstoneDev/pokedexcli/internal/pokecache"
)


func main() {
	// All available commands
	commands := functions.GetCommands()

	// create config
	currentConfig := functions.Config {
		Next: "",
		Previous : "",
	}

	// create cache
	apiCache := *pokecache.NewCache(5 * time.Minute)

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

		// get the command Name & the argument passed
		cmdName, arg := functions.NameAndArg(input)

		if cmd, ok := commands[cmdName]; ok {
			err := cmd.Callback(&currentConfig, &apiCache, arg)

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
