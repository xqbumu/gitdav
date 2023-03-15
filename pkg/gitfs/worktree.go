package gitfs

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"golang.org/x/net/webdav"
)

type FileSystem struct {
	*git.Worktree
}

func NewFilesystem(worktree *git.Worktree) webdav.FileSystem {
	return &FileSystem{Worktree: worktree}
}

func (fs *FileSystem) Mkdir(ctx context.Context, path string, mode os.FileMode) error {
	return fs.Filesystem.MkdirAll(path, mode)
}

func (fs *FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if flag&os.O_CREATE == os.O_CREATE {
		newFile, err := fs.Filesystem.Create(name)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()
	}

	stat, err := fs.Filesystem.Stat(name)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return NewDir(fs.Worktree, name), nil
	}

	return NewFile(fs.Worktree, name, flag, perm), nil
}

func (fs *FileSystem) RemoveAll(ctx context.Context, name string) error {
	return fs.Filesystem.Remove(name)
}

func (fs *FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return fs.Filesystem.Rename(oldName, newName)
}

func (fs *FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return fs.Filesystem.Stat(name)
}
