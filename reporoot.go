// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// repoRoot extracts the repository root path (host/owner/repo) from a URL
// string. For well-known hosts, the root is always the first 3 path segments.
// The returned root is suitable for use as a local directory path.
func repoRoot(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}

	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("no host in URL %q", rawURL)
	}

	cleaned := strings.TrimPrefix(path.Clean(u.Path), "/")
	cleaned = strings.TrimSuffix(cleaned, ".git")

	parts := strings.Split(cleaned, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("URL %q does not contain owner/repo", rawURL)
	}

	// Use only owner/repo (first two path segments) to form the root.
	return path.Join(host, parts[0], parts[1]), nil
}
