package main

type userRepoLanguages struct {
	User         string         `json:"username"`
	Repositories []Repositories `json:"repositories"`
}

type Repositories struct {
	Name      string   `json:"repositoryName"`
	Languages []string `json:"languages"`
}
