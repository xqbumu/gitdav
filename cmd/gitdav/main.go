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

	flag.Parse()
	if len(flag.Args()) != 1 || *commit == "" {
		flag.Usage()
		os.Exit(2)
	}

	repo := gitfs.NewRepo(flag.Args()[0])

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

	log.Println("serving requests for", repo.Cwd(), "at commit", repo.Commit())
	log.Fatalf("%+v", http.ListenAndServe(*httpAddr, &dav))
}
