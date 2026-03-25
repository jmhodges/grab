// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/url"
	"strings"
)

// rewriteToSSH rewrites an HTTPS repository URL to its SSH equivalent
// (git@host:path.git) if the URL's host is in sshHosts.
func rewriteToSSH(repoURL string, sshHosts []string) string {
	if len(sshHosts) == 0 {
		return repoURL
	}
	u, err := url.Parse(repoURL)
	if err != nil {
		return repoURL
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return repoURL
	}
	host := u.Hostname()
	for _, h := range sshHosts {
		if strings.EqualFold(host, h) {
			path := strings.TrimPrefix(u.Path, "/")
			if !strings.HasSuffix(path, ".git") {
				path += ".git"
			}
			return "git@" + host + ":" + path
		}
	}
	return repoURL
}
