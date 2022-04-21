package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

/*
Read through all of the users repositories, then iterate through the repositories and read all of the
languages github recognizes.
*/
func readRepositories(user string) (userRepoLanguages, error) {

	//Declare the necessary variables and prepare the HTTP Request
	userUrls := make([]string, 0)
	var userRepos userRepoLanguages
	client := &http.Client{}
	pages := 1
	req, _ := http.NewRequest("GET", "https://api.github.com/users/", nil)
	req.URL.Path += user + "/repos"
	q := req.URL.Query()

	//This is because of the pagination that github implemented in its API and should help with the 60 query limit per hour
	q.Add("page", "1")
	q.Add("per_page", "100")

	//Keep sending HTTP requests until we have all of the data
	for {
		//This variable is used to Unmarshal the API return
		var queryInfo []githubInfo

		//Finish setting up the HTTP Request, send it and read the body
		q.Set("page", strconv.Itoa(pages))
		req.URL.RawQuery = q.Encode()
		response, _ := client.Do(req)
		body, _ := io.ReadAll(response.Body)

		//Unmarshal the byte array into the []struct and close the body not that we're done with it
		json.Unmarshal(body, &queryInfo)
		response.Body.Close()

		//Iterate through the []Repos we received and append the LanguageURL for each user repository
		for _, githubRepo := range queryInfo {
			userUrls = append(userUrls, githubRepo.LanguagesURL)
		}

		//Return the current list of users and an error if we have hit the limit because we won't be able
		//to process anything further. This allows the program to write what it has rather than panic
		if response.Header["X-Ratelimit-Remaining"][0] == "0" {
			fmt.Println("Requests reset: ", response.Header["X-Ratelimit-Reset"])
			return userRepos, errors.New("Request limit hit, please try again later")
		}

		//Check if headers "Link" is there and if next is a part of it
		//Added the "next" check to keep you from burning all 60 requests and getting a 60 minute timeout.... Ask me how I know...
		if headerLinks, exists := response.Header["Link"]; !exists || !strings.Contains(headerLinks[0], "next") {
			break
		} else {
			pages++
		}

	}

	//Start building our struct
	userRepos.User = user

	//Iterate through the URLS we received and find the languages associated with the Repos
	for _, repoURL := range userUrls {
		//Declaring variables for scope
		//This is probably not the best way to do this but I find it much easier to read
		var repos Repositories
		var err error
		name := strings.Split(repoURL, "/")

		//Read the languages and if we hit the github rate limit (An err) return what we have so far
		if repos, err = readLanguages(repoURL); err != nil {
			return userRepos, err
		}

		//This is to pull the second to last from the array
		repos.Name = name[len(name)-2]

		//Append the repository information to the struct we started building
		userRepos.Repositories = append(userRepos.Repositories, repos)
	}

	return userRepos, nil
}

/*
Read the languages of a repository using the passed in url and return the Repositories struct and an error
The only error that should be returned is the github rate limit.
*/
func readLanguages(repository string) (Repositories, error) {
	//Declare the necessary variables and prepare the HTTP Request
	var repo Repositories
	client := &http.Client{}
	pages := 1
	req, _ := http.NewRequest("GET", repository, nil)
	q := req.URL.Query()

	//This *SHOULD* be unneccessary... I... I'm terrified of any project that has more than 100 languages...
	q.Add("page", "1")
	q.Add("per_page", "100")

	//Keep sending HTTP requests until we have all of the data
	for {
		//Declare an interface for Unmarshalling
		var lang map[string]interface{}

		//Finish setting up the HTTP Request, send it and read the body
		q.Set("page", strconv.Itoa(pages))
		req.URL.RawQuery = q.Encode()
		response, _ := client.Do(req)
		body, _ := io.ReadAll(response.Body)

		//This is a map, each key is the language and the value is a numerical representation
		json.Unmarshal(body, &lang)

		//Iterate through the keys of the map so we can append it to our Repositories struct
		for key := range lang {
			repo.Languages = append(repo.Languages, key)
		}

		response.Body.Close()

		//Return the current list of Repositories and an error if we have hit the limit because we won't be able
		//to process anything further. This allows the program to write what it has rather than panic
		if response.Header["X-Ratelimit-Remaining"][0] == "0" {
			fmt.Println("Requests reset: ", response.Header["X-Ratelimit-Reset"])
			return repo, errors.New("Request limit hit, please try again later")
		}

		//Check if headers "Link" is there and if next is a part of it
		if headerLinks, exists := response.Header["Link"]; !exists || !strings.Contains(headerLinks[0], "next") {
			break
		} else {
			pages++
		}
	}

	return repo, nil
}
