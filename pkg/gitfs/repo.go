package gitfs

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/net/webdav"
)

type Repo struct {
	cwd          string
	repo         *git.Repository
	checkoutOpts *git.CheckoutOptions
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

	return &Repo{cwd: path, repo: repo, checkoutOpts: &git.CheckoutOptions{}}
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
