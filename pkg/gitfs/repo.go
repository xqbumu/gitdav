package gitfs

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
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

	head, err := repo.Head()
	if err != nil {
		panic(err)
	}

	return &Repo{
		location: location,
		cwd:      path,
		repo:     repo,
		checkoutOpts: &git.CheckoutOptions{
			Branch: plumbing.ReferenceName(head.Name().String()),
		},
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

func (r *Repo) Head() string {
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

func (r *Repo) GetFileSystem() webdav.FileSystem {
	worktree, err := r.repo.Worktree()
	if err != nil {
		panic(err)
	}

	err = worktree.Checkout(r.checkoutOpts)
	if err != nil {
		panic(err)
	}

	return NewFilesystem(worktree)
}

func (r *Repo) Walk(fn filepath.WalkFunc) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	walkFileInfo("/", worktree.Filesystem, fn)

	return nil
}

func walkFileInfo(root string, fs billy.Filesystem, fn filepath.WalkFunc) {
	file, err := fs.Lstat(root)
	if err != nil {
		fn(root, nil, err)
	}
	if !file.IsDir() {
		fn(root, file, nil)
		return
	}
	files, err := fs.ReadDir(root)
	if err != nil {
		fn(root, nil, err)
	}
	for _, file := range files {
		walkFileInfo(filepath.Join(root, file.Name()), fs, fn)
	}
}
