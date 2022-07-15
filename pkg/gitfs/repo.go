package gitfs

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/net/webdav"
)

type Repo struct {
	location     string
	cwd          string
	repo         *git.Repository
	checkoutOpts *git.CheckoutOptions
}

func NewRepo(path string) *Repo {
	var err error
	location := "local"
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "git@") {
		location = "remote"
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			panic(err)
		}
	}

	var repo *git.Repository
	switch location {
	case "local":
		repo, err = git.PlainOpen(path)
	case "remote":
		repo, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{URL: path})
	}
	if err != nil {
		panic(err)
	}

	return &Repo{
		location:     location,
		cwd:          path,
		repo:         repo,
		checkoutOpts: &git.CheckoutOptions{},
	}
}

func (r *Repo) Branch(name string) (*config.Branch, error) {
	branch, err := r.repo.Branch(name)
	if err != nil {
		return nil, err
	}

	return branch, err
}

func (r *Repo) Cwd() string {
	return r.cwd
}

func (r *Repo) HEAD() string {
	ref, err := r.repo.Head()
	if err != nil {
		panic(err)
	}

	return ref.String()
}

func (r *Repo) SetCommit(commit string) {
	r.checkoutOpts = &git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	}
}

func (r *Repo) SetBranch(branch string, force bool) {
	r.checkoutOpts = &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  force,
	}
}

func (r *Repo) GetDir() webdav.FileSystem {
	worktree, err := r.repo.Worktree()
	if err != nil {
		panic(err)
	}

	worktree.Checkout(r.checkoutOpts)

	return NewFilesystem(worktree.Filesystem)
}
