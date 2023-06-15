package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v53/github"
)

func main() {
	ctx := context.Background()

	// TokenSource
	ts, err := Token()
	if err != nil {
		log.Println(err)
	}

	// TokenClient
	tc := TokenClient(ctx, ts)

	client := github.NewClient(tc)

	user := "luisedmc"

	// searching "dsa" repo in "luisedmc" user
	repo, _, err := client.Repositories.Get(ctx, user, "dsa")
	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}

	repoData := &Repository{
		FullName:    *repo.FullName,
		Description: *repo.Description,
		URL:         *repo.Owner.HTMLURL,
	}

	fmt.Println("Owner: ", repoData.FullName)
	fmt.Println("Repo Description: ", repoData.Description)
	fmt.Println("Repo URL: ", repoData.URL)
}
