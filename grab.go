// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	mvcs "github.com/Masterminds/vcs"
)

// grabDeps holds injectable dependencies for the grab function.
type grabDeps struct {
	newRepo       func(remote, local string) (mvcs.Repo, error)
	resolveImport func(importPath string) (goImport, error)
}

// grab resolves, fetches, and stores a repository for the given URL.
func grab(rawURL string, cfg config, deps grabDeps) error {
	// Ensure the input has a scheme so we can parse it as a URL.
	u, err := url.Parse(rawURL)
	if err == nil && u.Scheme == "" {
		rawURL = "https://" + rawURL
	}

	// Strip the scheme to get the import path for go-import resolution.
	importPath := strings.TrimPrefix(strings.TrimPrefix(rawURL, "https://"), "http://")

	// Try direct repo root derivation first. If that fails (e.g. vanity
	// import like go.uber.org/zap), fall back to go-import meta tag
	// resolution.
	root, err := repoRoot(rawURL)
	remoteURL := rawURL
	if err != nil {
		gi, giErr := deps.resolveImport(importPath)
		if giErr != nil {
			return fmt.Errorf("unable to determine repo from %q: %w", rawURL, giErr)
		}
		root = gi.Root
		remoteURL = gi.RepoURL
	}

	remoteURL = rewriteToSSH(remoteURL, cfg.SSHPreferredHosts)

	localPath := filepath.Join(cfg.Home, root)
	repo, err := deps.newRepo(remoteURL, localPath)
	if err != nil {
		return fmt.Errorf("unable to determine VCS type for %q: %w", remoteURL, err)
	}
	err = repo.Get()
	if err != nil {
		return fmt.Errorf("unable to download %q into %q: %w", remoteURL, localPath, err)
	}
	return nil
}
