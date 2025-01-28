package main

import (
	"fmt"
    "os"
    "context"
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
    "github.com/go-git/go-git/v5/plumbing/transport/http"
    "github.com/pkg/errors"
    "github.com/go-git/go-git/v5/plumbing/transport/client"
    "github.com/hashicorp/go-cleanhttp"
)


type Git struct {
	auth *http.BasicAuth
}

func NewGitWithAuth(auth *http.BasicAuth) *Git {
	client.InstallProtocol(
		"https",
		http.NewClient(cleanhttp.DefaultPooledClient()),
	)

	return &Git{auth: auth}
}

func NewGitClient() *Git {
	return NewGitWithAuth(&http.BasicAuth{
		Username: "user",
		Password: os.Getenv("GITLAB_TOKEN"),
	})
}

func (g *Git) CloneGitRepo(
	ctx context.Context,
	repoURL, branch, path string,
	depth int,
) (*git.Repository, error) {
	repository, err := git.PlainCloneContext(
		ctx,
		path,
		false,
		&git.CloneOptions{
			URL:               repoURL,
			Auth:              g.auth,
			ReferenceName:     plumbing.NewBranchReferenceName(branch),
			SingleBranch:      true,
			Depth:             depth,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			ShallowSubmodules: true,
			Tags:              git.NoTags,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error cloning Git repository")
	}

	return repository, nil
}

func main() {
	ctx := context.Background()

    repoURL := "https://github.com/aanogueira/go-git-issue.git"
    clonedBranch := "master"
    projectPath := "~/testerino"
    depth := 1

	client := NewGitClient()
	_, err := client.CloneGitRepo(
		ctx,
		repoURL,
		clonedBranch,
		projectPath,
		depth,
	)
	if err != nil {
        os.Exit(1)
	}

    return 
}
