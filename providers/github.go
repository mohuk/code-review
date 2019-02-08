package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/v22/github"
	gittypes "github.com/mohuk/code-review/types"
	"golang.org/x/oauth2"
)

var ctx context.Context
var location = "Asia/Karachi"

// GithubReview struct
type GithubReview struct {
	AccessToken string
	Client      *github.Client
	Repository  string
	UserAccount string
}

// GithubCreateBranchResponse struct
type GithubCreateBranchResponse struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

// GithubCreateBranchRequest struct
type GithubCreateBranchRequest struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

// InitializeClient creates a client for Github API
func (gr *GithubReview) InitializeClient() (*github.Client, error) {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gr.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	gr.Client = github.NewClient(tc)
	return gr.Client, nil
}

// CreateReviewPullRequest creates a PR for review
func (gr *GithubReview) CreateReviewPullRequest(base string, head string) (*github.PullRequest, error) {
	loc, _ := time.LoadLocation(location)
	today := time.Now().In(loc)
	prefix := fmt.Sprintf("%d-%d-%d", today.Day(), today.Month(), today.Year())
	title := fmt.Sprintf("Code Review (%s)", prefix)
	opts := github.NewPullRequest{
		Base:  &base,
		Title: &title,
		Head:  &head,
	}

	pr, _, err := gr.Client.PullRequests.Create(ctx, gr.UserAccount, gr.Repository, &opts)

	if err != nil {
		return nil, err
	}

	return pr, nil
}

// CreateReviewBranch creates a branch for code review
func (gr *GithubReview) CreateReviewBranch(base, name string) (*gittypes.GitBranch, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", gr.UserAccount, gr.Repository)

	payload := GithubCreateBranchRequest{}
	payload.Ref = fmt.Sprintf("refs/heads/%s", name)
	payload.Sha = gr.getLastCommit(base).Sha
	byt, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(byt))
	token := base64.StdEncoding.EncodeToString([]byte(gr.AccessToken))
	req.Header.Set("Authorization", "Basic "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	byts, _ := ioutil.ReadAll(res.Body)

	resBody := GithubCreateBranchResponse{}

	_ = json.Unmarshal(byts, &resBody)

	gitBranch := gittypes.GitBranch{}
	gitBranch.Sha = resBody.Object.Sha

	return &gitBranch, err
}

// BranchExists checks for existing branch
func (gr *GithubReview) BranchExists(branch string) (bool, error) {
	br, response, err := gr.Client.Repositories.GetBranch(ctx, gr.UserAccount, gr.Repository, branch)

	if response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	exists := branch == br.GetName()

	return exists, nil
}

func (gr *GithubReview) getLastCommit(branch string) (*gittypes.GitCommit, error) {
	commits, _, err := gr.Client.Repositories.ListCommits(ctx, gr.UserAccount, gr.Repository, nil)

	if err != nil {
		return nil, err
	}

	return &gittypes.GitCommit{
		Sha: *commits[len(commits)-1].SHA,
	}, nil
}
