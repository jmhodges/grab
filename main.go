// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/vcs"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) <= 1 || len(os.Args) > 2 {
		log.Fatal("usage: grab REPO_URL")
	}
	repoRoot, err := vcs.RepoRootForImportPath(os.Args[1], false)
	if err != nil {
		log.Fatal(err)
	}
	home := strings.TrimSpace(os.Getenv("GRAB_HOME"))
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("unable to get current user's home directory: %s", err)
		}
		home = filepath.Join(userHome, "src")
	}
	root := filepath.Join(home, repoRoot.Root)
	err = repoRoot.VCS.Create(root, repoRoot.Repo)
	if err != nil {
		log.Fatalf("unable to download %#v into %#v", repoRoot.Repo, root)
	}
}
