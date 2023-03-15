package gitfs

import (
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"golang.org/x/net/webdav"
)

type File struct {
	*git.Worktree
	name   string
	flag   int
	perm   os.FileMode
	offset int64
}

type fileStat interface {
	Stat() (os.FileInfo, error)
}

func NewFile(worktree *git.Worktree, name string, flag int, perm os.FileMode) webdav.File {
	return &File{Worktree: worktree, name: name, flag: flag, perm: perm}
}

func (f *File) Close() error { return nil }
func (f *File) Read(p []byte) (int, error) {
	file, err := f.Filesystem.OpenFile(f.name, f.flag, f.perm)
	if err != nil {
		return 0, err
	}
	return file.Read(p)
}

func (f *File) Readdir(int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	switch {
	case whence == io.SeekStart:
		f.offset = 0
	case whence == io.SeekEnd:
		stat, err := f.Stat()
		if err != nil {
			return 0, err
		}
		f.offset = stat.Size()
	default:
		return 0, os.ErrInvalid
	}

	f.offset += offset
	return f.offset, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.Filesystem.Stat(f.name)
}

func (f *File) Write(p []byte) (int, error) {
	file, err := f.Filesystem.OpenFile(f.name, f.flag, f.perm)
	if err != nil {
		return 0, err
	}

	n, err := file.Write(p)
	if err != nil {
		return 0, err
	}

	_, err = f.Add(f.name)
	if err != nil {
		return 0, err
	}

	return n, err
}
