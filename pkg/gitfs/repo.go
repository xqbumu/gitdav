package gitfs

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"golang.org/x/net/webdav"
)

type Repo struct {
	cwd    string
	commit string
	repo   *git.Repository
}

func NewRepo(path string) *Repo {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	return &Repo{cwd: path, repo: repo}
}

func (r *Repo) Cwd() string {
	return r.cwd
}

func (r *Repo) Commit() string {
	return r.commit
}

func (r *Repo) GetDir() webdav.FileSystem {
	worktree, err := r.repo.Worktree()
	if err != nil {
		panic(err)
	}

	return NewFilesystem(worktree.Filesystem)
}
