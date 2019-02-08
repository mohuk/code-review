# code-review

A package which provides an API which encapsulates a reasonable recurring code review process.

Steps:

1) Create a new branch say `code-review` from default branch say `master`
2) Continue working on `master`
3) Create a Pull request from `master` to `code-review` say `pr#1`
4) Review code on source control and report issues
5) Merge pull request (`pr#1`) once review items have been addressed
6) Return to step 2)

Package provides a `cli` for the entire process (WIP)

## Build

```bash
# package uses Go Modules
$ GO111MODULE=on go build ./cmd/code-review/main.go
```

## Run
```bash
$ ./main -help
Usage of ./main:
  -branch string
    	Name for review branch (default "code-review")
  -head string
    	Head for code review pull request
  -provider string
    	Version control provider (default "github")
  -repo string
    	Repository (default "test-repo")
  -task string
    	Task to perform
    	 initiate
    	 close
  -token string
    	Access Token (default "loasiqjjhq18982jahsd8912")
  -user string
    	Username (default "mohuk")
```
