package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

/*
Read a given user file. This method does check if the file exists, is a .txt file,
and can open the file. If this program is unable to complete any of these, it exits
and notifies the user of the issue.
*/
func readFile(inputFile string) []string {

	//Create required variables, this allows for a much nicer error checking
	users := []string{}
	var file *os.File
	var err error

	//This makes the user use a .txt file.
	// ASSUMPTION
	if !strings.HasSuffix(inputFile, ".txt") {
		log.Fatal("Input File is not a .txt file.")
	}

	//If file doesn't exist, we have nothing to do. Exit and notify user.
	if _, err = os.Stat(inputFile); err != nil {
		log.Fatal("Input File doesn't exist!")
	}

	//Attempt to open the file -- Return an empty array and notify via error that we we
	if file, err = os.Open(inputFile); err != nil {
		log.Fatal("Unable to read the file")
	}

	//Create a scanner, this could run into memory limitations
	// ASSUMPTION
	scanner := bufio.NewScanner(file)

	//Iterate through the file, ignoring empty lines.
	for scanner.Scan() {
		user := scanner.Text()

		//Check this to see if it's an empty line
		// ASSUMPTION
		if len(user) > 0 {
			users = append(users, user)
		}
	}

	//Debated on returning an error but since this is a critical step and we are just stopping execution
	//This would be only for future expansions.
	return users
}

/*
Write to the user given file. This file does check if it can open the file (in create,
truncate and WR), and can Marshal the json data. If either of these error it exits the
program and notifies the user. If the
*/
func writeFile(outputFile string, jsonData []userRepoLanguages) error {

	//Declaring these up front makes for cleaner error handling
	var fileHandler *os.File
	var marshalledData []byte
	var err error

	//Using MarshalIndent so it's easier to read. It consumes more space but it's actually ledgible
	//If an error occurs immediately notify the user and terminate as we have nothing to do
	if marshalledData, err = json.MarshalIndent(jsonData, "", "	"); err != nil {
		log.Fatal("Issues Marshalling the json Data")
	}

	//Warn the user they didn't give a .json file and append it to the given filename
	// ASSUMPTION
	if !strings.HasSuffix(outputFile, ".json") {
		log.Println("Output file doesn't have the proper file format. Appending .json...")
		outputFile = outputFile + ".json"
	}

	//Open the file WR, truncating the contents, creating it if it doesn't exist and in (rw- --r --r).
	//If an error occurs immediately notify the user and terminate as we have nothing to do
	if fileHandler, err = os.OpenFile(outputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 644); err != nil {
		log.Fatal("Issues opening the output file")
	}

	//Actually write to file and if an error occurs close the file and immediately notify the user and exit
	if _, err = fileHandler.Write(marshalledData); err != nil {
		err = fileHandler.Close()
		log.Fatal("Issues writing to file: ", err)
	}

	//This will either return nil (Successful) or an error which we return next line either way
	err = fileHandler.Close()

	return err
}
