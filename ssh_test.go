// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import "testing"

func TestRewriteToSSH(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		sshHosts []string
		want     string
	}{
		{
			name:     "github https to ssh",
			repoURL:  "https://github.com/jmhodges/grab",
			sshHosts: []string{"github.com"},
			want:     "git@github.com:jmhodges/grab.git",
		},
		{
			name:     "github https with .git suffix",
			repoURL:  "https://github.com/jmhodges/grab.git",
			sshHosts: []string{"github.com"},
			want:     "git@github.com:jmhodges/grab.git",
		},
		{
			name:     "non-matching host unchanged",
			repoURL:  "https://gitlab.com/user/repo",
			sshHosts: []string{"github.com"},
			want:     "https://gitlab.com/user/repo",
		},
		{
			name:     "empty ssh hosts unchanged",
			repoURL:  "https://github.com/jmhodges/grab",
			sshHosts: nil,
			want:     "https://github.com/jmhodges/grab",
		},
		{
			name:     "non-https scheme unchanged",
			repoURL:  "git://github.com/jmhodges/grab",
			sshHosts: []string{"github.com"},
			want:     "git://github.com/jmhodges/grab",
		},
		{
			name:     "http also rewritten",
			repoURL:  "http://github.com/jmhodges/grab",
			sshHosts: []string{"github.com"},
			want:     "git@github.com:jmhodges/grab.git",
		},
		{
			name:     "multiple ssh hosts",
			repoURL:  "https://gitlab.com/user/repo",
			sshHosts: []string{"github.com", "gitlab.com"},
			want:     "git@gitlab.com:user/repo.git",
		},
		{
			name:     "case insensitive host match",
			repoURL:  "https://GitHub.com/jmhodges/grab",
			sshHosts: []string{"github.com"},
			want:     "git@GitHub.com:jmhodges/grab.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteToSSH(tt.repoURL, tt.sshHosts)
			if got != tt.want {
				t.Errorf("rewriteToSSH(%q, %v) = %q, want %q", tt.repoURL, tt.sshHosts, got, tt.want)
			}
		})
	}
}
