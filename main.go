package main

import (
	"fmt"
	"log"
)

func main() {
	// TokenSource
	ts, err := Token()
	if err != nil {
		log.Println(err)
	}

	// TokenClient
	tc := TokenClient(ts)

	client := GithubClient(tc)

	var (
		user     string
		repoName string
	)
	fmt.Scanf("%s", &user)
	fmt.Scanf("%s", &repoName)

	// searching "repoName" repo in "user"
	repo := SearchRepository(client, user, repoName)

	fmt.Println("Owner: ", repo.FullName)
	fmt.Println("Repo Description: ", repo.Description)
	fmt.Println("Repo URL: ", repo.URL)
}
