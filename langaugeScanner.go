package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {

	//Declaring variables and flags to deal with cli arguments.
	var usersRepos []userRepoLanguages

	inputFile := flag.String("input", "Default", "/Example/Path/To/Input/File.txt")
	outputFile := flag.String("output", "./output.json", "/Example/Path/To/Output/File.json")

	flag.Parse()

	//This makes the user use a .txt file. This might be a bit strict this is an assumption.
	if !strings.HasSuffix(*inputFile, ".txt") {
		log.Fatal("Input File is not a .txt file.")
	}

	//If file doesn't exist, we have nothing to do. Exit and notify user.
	if _, err := os.Stat(*inputFile); err != nil {
		log.Fatal("Input File doesn't exist!")
	}

	//Read the user file (ASSUMPTION HERE!)
	users, err := readFile(*inputFile)

	//Make sure we read the file appropriately. If something goes wrong, we can't continue
	if err != nil {
		log.Fatal("Error reading the file: ", err)
	}

	//Warn the user they didn't give a .json file and append it to the given filename
	if !strings.HasSuffix(*outputFile, ".json") {
		log.Println("Output file doesn't have the proper file format. Appending .json...")
		*outputFile = *outputFile + ".json"
	}

	//Iterate through the list of users.
	for _, user := range users {
		userRepo, err := readRepositories(user)
		usersRepos = append(usersRepos, userRepo)

		//Warn the user if there are any errors, exit the loop and then attempt to write
		//what was gathered to the file.
		if err != nil {
			log.Println(err)
			break
		}
	}

	//Output the gathered data to the file
	writeFile(*outputFile, usersRepos)

}
