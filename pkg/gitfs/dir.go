package gitfs

import (
	"os"

	"github.com/go-git/go-git/v5"
	"golang.org/x/net/webdav"
)

type Dir struct {
	*git.Worktree
	name string
}

func NewDir(worktree *git.Worktree, name string) webdav.File {
	return &Dir{Worktree: worktree, name: name}
}

func (d *Dir) Close() error               { return nil }
func (d *Dir) Read(p []byte) (int, error) { return 0, os.ErrInvalid }
func (d *Dir) Readdir(int) ([]os.FileInfo, error) {
	entries, err := d.Filesystem.ReadDir(d.name)
	if err != nil {
		return nil, err
	}

	items := []os.FileInfo{}
	for _, item := range entries {
		// ignore .git
		if item.Name() == ".git" {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (d *Dir) Stat() (os.FileInfo, error) {
	return d.Filesystem.Stat(d.name)
}

func (d *Dir) Write(p []byte) (int, error) {
	return 0, os.ErrInvalid
}
