package mod_git

import (
	"github.com/go-git/go-git/v5"
)

type GitDB struct {
	Path          string
	Repository    *git.Repository
	Worktree      *git.Worktree
	PullOptions   *git.PullOptions
	CommitOptions *git.CommitOptions
	PushOptions   *git.PushOptions
}
