// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// goImport holds the parsed contents of a <meta name="go-import"> tag.
type goImport struct {
	// Root is the import path prefix (e.g. "go.uber.org/zap").
	Root string
	// VCS is the version control system (e.g. "git").
	VCS string
	// RepoURL is the repository URL (e.g. "https://github.com/uber-go/zap").
	RepoURL string
}

// resolveGoImport fetches https://<importPath>?go-get=1 and parses the
// <meta name="go-import"> tag to discover the VCS type and repository URL.
func resolveGoImport(importPath string) (goImport, error) {
	u := "https://" + importPath + "?go-get=1"
	resp, err := http.Get(u)
	if err != nil {
		return goImport{}, fmt.Errorf("fetching %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return goImport{}, fmt.Errorf("fetching %s: status %d", u, resp.StatusCode)
	}

	return parseGoImport(resp.Body, importPath)
}

// parseGoImport reads HTML from r and returns the go-import metadata whose
// root is a prefix of importPath.
func parseGoImport(r io.Reader, importPath string) (goImport, error) {
	tokenizer := html.NewTokenizer(r)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return goImport{}, fmt.Errorf("no go-import meta tag found for %q", importPath)
		case html.SelfClosingTagToken, html.StartTagToken:
			tn, hasAttr := tokenizer.TagName()
			if string(tn) != "meta" || !hasAttr {
				continue
			}
			var name, content string
			for {
				key, val, more := tokenizer.TagAttr()
				switch string(key) {
				case "name":
					name = string(val)
				case "content":
					content = string(val)
				}
				if !more {
					break
				}
			}
			if name != "go-import" || content == "" {
				continue
			}
			gi, err := parseGoImportContent(content)
			if err != nil {
				continue
			}
			if strings.HasPrefix(importPath, gi.Root) {
				return gi, nil
			}
		}
	}
}

// parseGoImportContent parses the content attribute value of a go-import meta
// tag. The format is "<root> <vcs> <repo-url>".
func parseGoImportContent(content string) (goImport, error) {
	fields := strings.Fields(content)
	if len(fields) != 3 {
		return goImport{}, fmt.Errorf("invalid go-import content: %q", content)
	}
	return goImport{
		Root:    fields[0],
		VCS:     fields[1],
		RepoURL: fields[2],
	}, nil
}
