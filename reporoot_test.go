// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import "testing"

func TestRepoRoot(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "https github url",
			rawURL: "https://github.com/jmhodges/grab",
			want:   "github.com/jmhodges/grab",
		},
		{
			name:   "https with .git suffix",
			rawURL: "https://github.com/jmhodges/grab.git",
			want:   "github.com/jmhodges/grab",
		},
		{
			name:   "https with subpath",
			rawURL: "https://github.com/jmhodges/grab/subpkg/deep",
			want:   "github.com/jmhodges/grab",
		},
		{
			name:   "https with trailing slash",
			rawURL: "https://github.com/jmhodges/grab/",
			want:   "github.com/jmhodges/grab",
		},
		{
			name:   "double slashes in path",
			rawURL: "https://github.com//jmhodges//grab",
			want:   "github.com/jmhodges/grab",
		},
		{
			name:   "gitlab url",
			rawURL: "https://gitlab.com/user/repo",
			want:   "gitlab.com/user/repo",
		},
		{
			name:    "no path segments",
			rawURL:  "https://github.com",
			wantErr: true,
		},
		{
			name:    "only one path segment",
			rawURL:  "https://github.com/jmhodges",
			wantErr: true,
		},
		{
			name:    "no host",
			rawURL:  "/jmhodges/grab",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repoRoot(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("repoRoot(%q) = %q, want error", tt.rawURL, got)
				}
				return
			}
			if err != nil {
				t.Errorf("repoRoot(%q) unexpected error: %s", tt.rawURL, err)
				return
			}
			if got != tt.want {
				t.Errorf("repoRoot(%q) = %q, want %q", tt.rawURL, got, tt.want)
			}
		})
	}
}
