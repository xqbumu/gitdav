package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/xqbumu/gitdav/pkg/gitfs"
	"golang.org/x/exp/slog"
	"golang.org/x/net/webdav"
)

const (
	defaultAddr = ":6060" // default webserver address
)

func main() {
	httpAddr := flag.String("http", defaultAddr, "HTTP service address (e.g., '"+defaultAddr+"')")
	commit := flag.String("commit", "", "commit hash")
	branch := flag.String("branch", "", "working branch")
	create := flag.Bool("create", false, "create branch if branch not exist")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	repo := gitfs.NewRepo(flag.Args()[0])
	if len(*commit) > 0 {
		repo.SetCommit(*commit)
	} else if len(*branch) > 0 {
		repo.SetBranch(*branch, *create)
	}

	davPrefix := "/webdav"
	davHandler := NewDavHandler(davPrefix, repo)

	mux := chi.NewMux()
	mux.Mount(davPrefix, davHandler)
	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Debug\n")
		fmt.Fprintf(w, "Cwd:\t%s\n", repo.Cwd())
		fmt.Fprintf(w, "Head:\t%s\n", repo.Head())
		buf := bytes.NewBuffer(nil)
		err := repo.Walk(func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			if info.IsDir() {
				fmt.Fprintf(buf, "%s/%s/\n", path, info.Name())
			} else {
				fmt.Fprintf(buf, "%s/%s\n", path, info.Name())
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Files:\n%s\n", buf.String())
	})

	slog.Info(fmt.Sprintf("serving requests for %s, at %s", repo.Cwd(), repo.Head()))
	slog.With("err", http.ListenAndServe(*httpAddr, mux)).Warn("shutdown")
}

func NewDavHandler(prefix string, repo *gitfs.Repo) http.Handler {
	return &webdav.Handler{
		Prefix:     prefix,
		FileSystem: repo.GetFileSystem(),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			slog.
				With("proto", req.Proto).
				With("method", req.Method).
				With("url", req.URL).
				Info("request")
			if err != nil {
				slog.With("err", err).Error("error")
				return
			}
		},
	}
}
