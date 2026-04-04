// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"
)

func TestParseGoImport(t *testing.T) {
	tests := []struct {
		name       string
		html       string
		importPath string
		want       goImport
		wantErr    bool
	}{
		{
			name: "standard go-import tag",
			html: `<html><head>
				<meta name="go-import" content="go.uber.org/zap git https://github.com/uber-go/zap">
			</head></html>`,
			importPath: "go.uber.org/zap",
			want: goImport{
				Root:    "go.uber.org/zap",
				VCS:     "git",
				RepoURL: "https://github.com/uber-go/zap",
			},
		},
		{
			name: "import path is subpackage of root",
			html: `<html><head>
				<meta name="go-import" content="go.uber.org/zap git https://github.com/uber-go/zap">
			</head></html>`,
			importPath: "go.uber.org/zap/zapcore",
			want: goImport{
				Root:    "go.uber.org/zap",
				VCS:     "git",
				RepoURL: "https://github.com/uber-go/zap",
			},
		},
		{
			name: "self-closing meta tag",
			html: `<html><head>
				<meta name="go-import" content="gopkg.in/yaml.v2 git https://gopkg.in/yaml.v2" />
			</head></html>`,
			importPath: "gopkg.in/yaml.v2",
			want: goImport{
				Root:    "gopkg.in/yaml.v2",
				VCS:     "git",
				RepoURL: "https://gopkg.in/yaml.v2",
			},
		},
		{
			name: "multiple meta tags picks matching one",
			html: `<html><head>
				<meta name="go-import" content="example.com/other git https://github.com/other/repo">
				<meta name="go-import" content="example.com/foo git https://github.com/foo/repo">
			</head></html>`,
			importPath: "example.com/foo",
			want: goImport{
				Root:    "example.com/foo",
				VCS:     "git",
				RepoURL: "https://github.com/foo/repo",
			},
		},
		{
			name: "ignores non-go-import meta tags",
			html: `<html><head>
				<meta name="viewport" content="width=device-width">
				<meta name="go-import" content="example.com/foo git https://github.com/foo/repo">
			</head></html>`,
			importPath: "example.com/foo",
			want: goImport{
				Root:    "example.com/foo",
				VCS:     "git",
				RepoURL: "https://github.com/foo/repo",
			},
		},
		{
			name:       "no matching meta tag",
			html:       `<html><head><meta name="go-import" content="other.com/bar git https://github.com/bar/repo"></head></html>`,
			importPath: "example.com/foo",
			wantErr:    true,
		},
		{
			name:       "no meta tags at all",
			html:       `<html><head><title>Hello</title></head></html>`,
			importPath: "example.com/foo",
			wantErr:    true,
		},
		{
			name:       "malformed content attribute",
			html:       `<html><head><meta name="go-import" content="only-two-fields git"></head></html>`,
			importPath: "only-two-fields",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGoImport(strings.NewReader(tt.html), tt.importPath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseGoImport() = %+v, want error", got)
				}
				return
			}
			if err != nil {
				t.Errorf("parseGoImport() unexpected error: %s", err)
				return
			}
			if got != tt.want {
				t.Errorf("parseGoImport() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
