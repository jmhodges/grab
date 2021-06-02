package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/vcs"
)

func main() {
	repoRoot, err := vcs.RepoRootForImportPath(os.Args[1], true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", repoRoot)
	home := os.Getenv("GRAB_HOME")
	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("unable to get current user's home directory: %s", err)
		}
	}
	root := filepath.Join(home, "src", repoRoot.Root)
	fmt.Printf("root: %#v\n", root)
	err = repoRoot.VCS.Create(root, repoRoot.Repo)
	if err != nil {
		log.Fatalf("unable to download %#v into %#v: %#v", repoRoot.Repo, root, err)
	}
}
