// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	mvcs "github.com/Masterminds/vcs"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) <= 1 || len(os.Args) > 2 {
		log.Fatal("usage: grab REPO_URL")
	}
	rawURL := os.Args[1]

	// Ensure the input has a scheme so we can parse it as a URL.
	u, err := url.Parse(rawURL)
	if err == nil && u.Scheme == "" {
		rawURL = "https://" + rawURL
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err)
	}

	// Strip the scheme to get the import path for go-import resolution.
	importPath := strings.TrimPrefix(strings.TrimPrefix(rawURL, "https://"), "http://")

	// Try direct repo root derivation first. If that fails (e.g. vanity
	// import like go.uber.org/zap), fall back to go-import meta tag
	// resolution.
	root, err := repoRoot(rawURL)
	remoteURL := rawURL
	if err != nil {
		gi, giErr := resolveGoImport(importPath)
		if giErr != nil {
			log.Fatalf("unable to determine repo from %q: %s", rawURL, giErr)
		}
		root = gi.Root
		remoteURL = gi.RepoURL
	}

	remoteURL = rewriteToSSH(remoteURL, cfg.SSHPreferredHosts)

	localPath := filepath.Join(cfg.Home, root)
	repo, err := mvcs.NewRepo(remoteURL, localPath)
	if err != nil {
		log.Fatalf("unable to determine VCS type for %q: %s", remoteURL, err)
	}
	err = repo.Get()
	if err != nil {
		log.Fatalf("unable to download %#v into %#v: %s", remoteURL, localPath, err)
	}
}
