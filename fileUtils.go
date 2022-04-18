package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
)

func readFile(inputFile string) ([]string, error) {
	users := []string{}
	file, err := os.Open(inputFile)

	//Return an empty array and notify via error that we were unable to read the file
	if err != nil {
		return users, errors.New("Unable to read the file")
	}

	scanner := bufio.NewScanner(file)

	//Iterate through the file, ignoring empty lines.
	for scanner.Scan() {
		user := scanner.Text()

		//Check this to see if it's an empty line
		if len(user) > 0 {
			users = append(users, user)
		}
	}

	return users, nil
}

//Using unbuffered again. May not be the best idea...
func writeFile(outputFile string, jsonData []userRepoLanguages) error {

	//Declaring these up front makes for cleaner error handling
	var fileHandler *os.File
	var marshalledData []byte
	var err error

	//Open the file WR, truncating the contents, creating it if it doesn't exist and in (rw- --r --r).
	if fileHandler, err = os.OpenFile(outputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 644); err != nil {
		log.Println("Issues opening the output file")
		return err
	}

	//Using MarshalIndent so it's easier to read. It consumes more space but it's actually ledgible
	if marshalledData, err = json.MarshalIndent(jsonData, "", "	"); err != nil {
		log.Println("Issues Marshalling the json Data")
		return err
	}

	if _, err = fileHandler.Write(marshalledData); err != nil {
		log.Println("Issues writing to file: ", err)
		return err
	}

	//This will either return nil (Successful) or an error which we return next line either way
	err = fileHandler.Close()

	return err
}
