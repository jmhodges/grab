// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net/url"
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
	importPath := os.Args[1]

	// We parse as a URL to see if we need to strip out the leading scheme for
	// the vcs library. The vcs library works on "import paths" a la Go, not URLs.
	u, err := url.Parse(os.Args[1])
	// If we can't parse it as a URL, it might still mean the vcs library knows
	// how to handle it.
	if err == nil {
		u.Scheme = ""
		importPath = strings.TrimLeft(u.String(), "/")
	}

	repoRoot, err := vcs.RepoRootForImportPath(importPath, false)
	if err != nil {
		log.Fatalf("unable to figure out the repo root from the given url: %s", err)
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
