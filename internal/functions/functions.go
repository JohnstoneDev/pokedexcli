package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"

	"github.com/JohnstoneDev/pokedexcli/internal/pokecache"
)

// command type with name, description & a callback
type cliCommand struct {
	name        string
	description string
	Callback    func(configPtr *Config, cache *pokecache.Cache) error
}

type Config struct {
	Next string
	Previous string
}

// store all commands here
func GetCommands () map[string]cliCommand {
	return  map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"clear": {
			name: 			"clear",
			description: "Clears the console",
			Callback:		commandClear,
		},
		"map": {
			name: 			"map",
			description: "Displays the names of 20 location areas in the pokemon world.",
			Callback:		commandMap,
		},
		"mapb": {
			name: 			"mapb",
			description: "(Map Back) Displays the names of the last 20 location areas in the pokemon world you viewed.",
			Callback:		commandMapB,
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
func commandHelp (configPtr *Config, cache *pokecache.Cache) error {
	Welcome(GetCommands())
	return nil
}

// clear command : clears the console
func commandClear(configPtr *Config, cache *pokecache.Cache) error {
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
func commandExit (configPtr *Config, cache *pokecache.Cache) error {
	fmt.Println("Pokedex says bye bye!")
	return errors.New("")
}

// helper function that displays the data from the API
func displayItems(items []structs.Result) {
	for _, item := range items {
		fmt.Println(item.Name)
	}
}

// API Call to PokeAPI, plays with the config struct a bit
func commandMap (configPtr *Config, cache *pokecache.Cache) error {
	if configPtr.Next == "" {
		response, err := pokeapi.Resource("location")

		// set the previous link for the first request to the base for the API
		configPtr.Previous = "https://pokeapi.co/api/v2/location"
		configPtr.Next = response.Next

		if err != nil {
			fmt.Println(err)
			return err
		}

		// Add the response data to the cache
		cache.Add(configPtr.Previous, response)

		displayItems(response.Results)

	} else {
		// check if there is any information in the cache
		data, cached := cache.Get(configPtr.Next)

		// display the cached items or call the API
		if cached {
			displayItems(data.Results)

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

			// Add response to cache
			cache.Add(configPtr.Next, responseData)

			configPtr.Previous = configPtr.Next
			configPtr.Next = responseData.Next

			displayItems(responseData.Results)
		}
	}

	return nil
}

// Should show the previous 20 results
func commandMapB (configPtr *Config, cache *pokecache.Cache) error {
	if configPtr.Previous == "" {

		fmt.Println("Error: no previous request, please use map first")

	} else {
		// check if there are items in the cache
		data, cached := cache.Get(configPtr.Previous)

		// if there is none, call the api
		if cached {
			displayItems(data.Results)

		} else {
			// get the previous URL from the config
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

			displayItems(responseData.Results)
		}
	}
	return nil
}
