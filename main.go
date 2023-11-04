package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// command type with name, description & a callback
type cliCommand struct {
	name        string
	description string
	callback    func() error
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

// help callback
func commandHelp () error {
	Welcome(GetCommands())
	return nil
}

// clear command
func commandClear() error {
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
func commandExit () error {
	return errors.New("")
}

func main() {
	// All available commands
	commands := GetCommands()

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
			err := cmd.callback()

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
