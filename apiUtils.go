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

func readRepositories(user string) (userRepoLanguages, error) {
	//This is what I realistically want, a map which uses user as key and contains a series of string urls.
	userUrls := make([]string, 0)
	var userRepos userRepoLanguages
	client := &http.Client{}
	pages := 1

	req, _ := http.NewRequest("GET", "https://api.github.com/users/", nil)
	req.URL.Path += user + "/repos"
	q := req.URL.Query()
	q.Add("page", "1")
	//This is to help with the 60 query limit per hour (?) that github has in place
	//Would like to notify the user if we hit that limit
	q.Add("per_page", "100")

	for {
		var queryInfo []githubInfo
		q.Set("page", strconv.Itoa(pages))
		req.URL.RawQuery = q.Encode()
		response, _ := client.Do(req)
		body, _ := io.ReadAll(response.Body)

		//This is a list of maps will need to iterate through the list to get
		json.Unmarshal(body, &queryInfo)

		response.Body.Close()

		for _, githubRepo := range queryInfo {
			userUrls = append(userUrls, githubRepo.LanguagesURL)
		}

		//Return an error if we have hit the limit because we won't be able to process anything further
		//This allows the program to write what it has then return rather than panic
		if response.Header["X-Ratelimit-Remaining"][0] == "0" {
			fmt.Println("Remaining requests: ", response.Header["X-Ratelimit-Remaining"])
			fmt.Println("Requests reset: ", response.Header["X-Ratelimit-Reset"])
			return userRepos, errors.New("Request limit hit, please try again later")
		}

		//Check if headers "Link" is there and if next is a part of it
		//Keeps you from burning all 60 requests and getting a 60 minute timeout.... Ask me how I know...
		if headerLinks, exists := response.Header["Link"]; !exists || !strings.Contains(headerLinks[0], "next") {
			break
		} else {
			pages++
		}

	}

	userRepos.User = user

	//This is going to be a little ugly...
	for _, repoURL := range userUrls {
		repos, err := readLanguages(repoURL)
		if err != nil {
			return userRepos, err
		}

		userRepos.Repositories = append(userRepos.Repositories, repos)
	}

	return userRepos, nil
}

func readLanguages(repository string) (Repositories, error) {
	var repo Repositories
	client := &http.Client{}
	pages := 1

	req, _ := http.NewRequest("GET", repository, nil)
	q := req.URL.Query()
	//This *SHOULD* be unneccessary... I... I'm terrified of any project that has more than 100 languages...
	q.Add("page", "1")
	q.Add("per_page", "100")

	for {
		var lang map[string]interface{}
		q.Set("page", strconv.Itoa(pages))
		req.URL.RawQuery = q.Encode()
		response, _ := client.Do(req)
		body, _ := io.ReadAll(response.Body)

		//This is a list of maps will need to iterate through the list to get
		json.Unmarshal(body, &lang)
		fmt.Println(repository)
		for key := range lang {

			repo.Languages = append(repo.Languages, key)
		}

		response.Body.Close()

		//Return an error if we have hit the limit because we won't be able to process anything further
		//This allows the program to write what it has then return rather than panic
		if response.Header["X-Ratelimit-Remaining"][0] == "0" {
			/*currentTime := time.Now()
			currentTime.Sub(time.Unix(strconv.ParseInt(response.Header["X-Ratelimit-Reset"][0], 10, 64), 0))
			time.Parse(response.Header["X-Ratelimit-Reset"][0])*/
			fmt.Println("Remaining requests: ", response.Header["X-Ratelimit-Remaining"])
			fmt.Println("Requests reset: ", response.Header["X-Ratelimit-Reset"])
			return repo, errors.New("Request limit hit, please try again later")
		}

		//Check if headers "Link" is there and if next is a part of it
		//Keeps you from burning all 60 requests and getting a 40 minute timeout.... Ask me how I know...
		if headerLinks, exists := response.Header["Link"]; !exists || !strings.Contains(headerLinks[0], "next") {
			break
		} else {
			pages++
		}
	}

	return repo, nil
}
