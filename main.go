package main

import (
	"context"
	"fmt"
	"log"
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
	client := GithubClient(tc)

	var repoName string
	fmt.Scanf("%s", &repoName)

	// CreateRepository(ctx, repoName, true, client)
	SearchRepository(ctx, client, "luisedmc", repoName)
}
