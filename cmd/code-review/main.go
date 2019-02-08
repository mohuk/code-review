package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohuk/code-review/providers"
)

var (
	useraccount string
	repository  string
	accessToken string
	provider    string
	task        string
	base        string
	head        string
	branch      string
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&useraccount, "user", os.Getenv("GITHUB_USERACCOUNT"), "Username")
	flag.StringVar(&repository, "repo", os.Getenv("GITHUB_REPO"), "Repository")
	flag.StringVar(&accessToken, "token", os.Getenv("GITHUB_TOKEN"), "Access Token")
	flag.StringVar(&provider, "provider", "github", "Version control provider")
	flag.StringVar(&task, "task", "", "Task to perform \n initiate \n close")
	flag.StringVar(&head, "head", "", "Head for code review pull request")
	flag.StringVar(&branch, "branch", "code-review", "Name for review branch")

	flag.Parse()

	if useraccount == "" || repository == "" || accessToken == "" {
		log.Fatal("Missing Parameters")
	}

	if task == "" {
		log.Fatal("No task specified")
	}
	cfg := providers.ProviderConfig{
		Name:        provider,
		AccessToken: accessToken,
		UserAccount: useraccount,
		Repository:  repository,
	}
	provider := providers.NewProvider(cfg)
	provider.InitializeClient()

	switch task {
	case "initiate":
		{
			br, err := provider.BranchExists(branch)

			if err != nil {
				log.Fatal(err)
			}

			if !br {
				_, err := provider.CreateReviewBranch(head, branch)

				if err != nil {
					log.Fatal(err)
				}
			}
			pr, err := provider.CreateReviewPullRequest(branch, head)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Pull request for review created at %s", *pr.URL)
		}
	default:
		{
			fmt.Println("No/Incorrect option provided")
		}
	}

}
