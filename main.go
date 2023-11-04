package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"
)

// command type with name, description & a callback
type cliCommand struct {
	name        string
	description string
	callback    func(configPtr *config) error
}

type config struct {
	Next string
	Previous string
}

// store all commands here
func GetCommands () map[string]cliCommand {
	return  map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"clear": {
			name: 			"clear",
			description: "Clears the console",
			callback:		commandClear,
		},
		"map": {
			name: 			"map",
			description: "Displays the names of 20 location areas in the pokemon world.",
			callback:		commandMap,
		},
		"mapb": {
			name: 			"mapb",
			description: "(Map Back) Displays the names of the last 20 location areas in the pokemon world you viewed.",
			callback:		commandMapB,
		},
	}
}

// Welcome command that will also be used as the help command
func Welcome(commands map[string]cliCommand) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")

	fmt.Println()

	for _, command := range commands {
		fmt.Println(command.name, ": ", command.description)
	}

	fmt.Println()

	return nil
}

// help callback : shows info on the tool
func commandHelp (configPtr *config) error {
	Welcome(GetCommands())
	return nil
}

// clear command : clears the console
func commandClear(configPtr *config) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

// exit callback
func commandExit (configPtr *config) error {
	fmt.Println("Pokedex says bye bye!")
	return errors.New("")
}

// API Call to PokeAPI, plays with the config struct a bit
func commandMap (configPtr *config) error {

	if configPtr.Next == "" {
		response, err := pokeapi.Resource("location")
		configPtr.Next = response.Next

		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, item := range response.Results {
			fmt.Println(item.Name)
		}
	} else {
		// call the API using a get request with the next variable stored
		// in the config
		response, err := http.Get(configPtr.Next)

		if err != nil {
			return err
		}

		defer response.Body.Close()

		// custom data type that is used as the API response
		var responseData structs.Resource
		decoder := json.NewDecoder(response.Body)

		err = decoder.Decode(&responseData)

		if err != nil {
			return err
		}

		configPtr.Previous = configPtr.Next
		configPtr.Next = responseData.Next

		for _, item := range responseData.Results {
			fmt.Println(item.Name)
		}

	}

	return nil
}

// Should show the previous 20 results
func commandMapB (configPtr *config) error {

	if configPtr.Previous == "" {

		fmt.Println("Error: no previous request, please use map first")

	} else {
		response, err := http.Get(configPtr.Previous)

		if err != nil {
			return err
		}

		defer response.Body.Close()

		// custom data type that is used as the API response
		var responseData structs.Resource
		decoder := json.NewDecoder(response.Body)

		err = decoder.Decode(&responseData)

		if err != nil {
			return err
		}

		configPtr.Next = responseData.Next

		for _, item := range responseData.Results {
			fmt.Println(item.Name)
		}
	}
	return nil
}

func main() {
	// All available commands
	commands := GetCommands()

	currentConfig := config {
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
			err := cmd.callback(&currentConfig)

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
