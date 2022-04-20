package main

import (
	"flag"
	"log"
)

/*
Reads a list of usernames from a given .txt file and query github for the users repositories, and the language
associated with that repository. It then outputs the data into a readable JSON format into a given .json file.
*/
func main() {

	//Declaring variables and flags to deal with cli arguments.
	var usersRepos []userRepoLanguages

	inputFile := flag.String("input", "Default", "/Example/Path/To/Input/File.txt")
	outputFile := flag.String("output", "./output.json", "/Example/Path/To/Output/File.json")

	flag.Parse()

	//Read the user file
	//ASSUMPTION
	users := readFile(*inputFile)

	//Iterate through the list of users, read the repositories they have and the languages those
	//repositories use. There are a slew of things that can go wrong and we will stop if any user
	//returns an error. This is primarily due to the github api limiter.
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
