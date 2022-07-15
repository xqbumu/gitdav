package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/xqbumu/gitdav/pkg/gitfs"
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

	dav := webdav.Handler{
		FileSystem: repo.GetDir(),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			log.Printf("%v %v %v\n", req.Method, req.URL, req.Proto)
			if err != nil {
				log.Printf("%+v", err)
				return
			}
		},
	}

	log.Printf("serving requests for %s, at %s", repo.Cwd(), repo.HEAD())
	log.Fatalf("%+v", http.ListenAndServe(*httpAddr, &dav))
}
