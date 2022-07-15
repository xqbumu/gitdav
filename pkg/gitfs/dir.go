package gitfs

import (
	"os"

	"github.com/go-git/go-billy/v5"
	"golang.org/x/net/webdav"
)

type Dir struct {
	root billy.Filesystem
	name string
}

func NewDir(root billy.Filesystem, name string) webdav.File {
	return &Dir{root: root, name: name}
}

func (d *Dir) Close() error               { return nil }
func (d *Dir) Read(p []byte) (int, error) { return 0, os.ErrInvalid }
func (d *Dir) Readdir(int) ([]os.FileInfo, error) {
	entries, err := d.root.ReadDir(d.name)
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
	return d.root.Stat(d.name)
}

func (d *Dir) Write(p []byte) (int, error) {
	return 0, os.ErrInvalid
}
