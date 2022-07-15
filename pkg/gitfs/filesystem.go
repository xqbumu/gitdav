package gitfs

import (
	"context"
	"os"

	"github.com/go-git/go-billy/v5"
	"golang.org/x/net/webdav"
)

type Filesystem struct {
	root billy.Filesystem
}

func NewFilesystem(root billy.Filesystem) webdav.FileSystem {
	return &Filesystem{root: root}
}

func (fs *Filesystem) Mkdir(ctx context.Context, path string, mode os.FileMode) error {
	return fs.root.MkdirAll(path, mode)
}

func (fs *Filesystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if flag&os.O_CREATE == os.O_CREATE {
		newFile, err := fs.root.Create(name)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()
	}

	stat, err := fs.root.Stat(name)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return NewDir(fs.root, name), nil
	}

	return NewFile(fs.root, name, flag, perm), nil
}

func (fs *Filesystem) RemoveAll(ctx context.Context, name string) error {
	return fs.root.Remove(name)
}

func (fs *Filesystem) Rename(ctx context.Context, oldName, newName string) error {
	return fs.root.Rename(oldName, newName)
}

func (fs *Filesystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return fs.root.Stat(name)
}
