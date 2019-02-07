package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	providers "github.com/mohuk/go-github-issue-report/providers"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	gReview := providers.GithubReview{
		AccessToken: os.Getenv("GITHUB_TOKEN"),
		UserAccount: os.Getenv("GITHUB_USERACCOUNT"),
		Repository:  os.Getenv("GITHUB_REPO"),
	}

	gReview.InitializeClient()
	sha, err := gReview.CreateReviewBranch("weekly-review", "92019d7fe70b8a54dafb69b902513f864f7c3f4f")

	if err != nil {
		panic(err)
	}

	fmt.Printf("Branch created from %s", *sha)

}
